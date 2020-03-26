package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

var c *http.Client = http.DefaultClient

func Init(proxy string) error {
	c = http.DefaultClient
	proxy = strings.TrimSpace(proxy)
	if proxy != "" {
		url, err := url.Parse(proxy)
		if err != nil {
			return err
		}

		c = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(url),
			},
		}
	}

	return nil
}

func Get(url string, headers map[string]string) ([]byte, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		r.Header.Add(k, v)
	}

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") || strings.HasPrefix(contentType, "video/") {
		return b, nil
	}

	encoding, _, _ := charset.DetermineEncoding(b, contentType)
	utf8Reader := transform.NewReader(bytes.NewReader(b), encoding.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func Post(url string, body map[string]interface{}) ([]byte, error) {
	bodyBytes := []byte("")
	if len(body) > 0 {
		var err error
		if bodyBytes, err = json.Marshal(body); err != nil {
			return nil, err
		}
	}

	r, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
