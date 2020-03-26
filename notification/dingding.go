package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Dingding struct {
	bot  string
	freq time.Duration
	last time.Time
	mu   sync.Mutex
}

func NewDingding(bot string, freq time.Duration) *Dingding {
	return &Dingding{
		bot:  bot,
		freq: freq,
	}
}

type Response struct {
	Msg  string `json:"errmsg"`
	Code int    `json:"errcode"`
}

func (d *Dingding) SendRated(title string, rows []string) error {
	d.mu.Lock()
	if time.Since(d.last) < d.freq {
		d.mu.Unlock()
		return nil
	}

	d.last = time.Now()
	d.mu.Unlock()

	return d.Send(title, rows)
}

func (d *Dingding) Send(title string, rows []string) error {
	fmtRows := []string{fmt.Sprintf("##### **%v**", title)}
	for i := 0; i < len(rows); i++ {
		fmtRows = append(fmtRows, fmt.Sprintf("###### %v", rows[i]))
	}

	msg := []byte(fmt.Sprintf(`{
"msgtype": "markdown",
"markdown": {
	"title": "%v",
	"text": "%v \n\n"
},
"at": {
	"atMobiles": [],
	"isAtAll": true
}
}`, title, strings.Join(fmtRows, " \n\n ")))

	req, err := http.NewRequest("POST", string(d.bot), bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("non-200 code received: %v", resp.StatusCode))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res Response
	if err := json.Unmarshal(bytes, &res); err != nil {
		return err
	}
	if res.Code != 0 {
		return errors.New(fmt.Sprintf("error response: %v", res.Msg))
	}
	return nil
}
