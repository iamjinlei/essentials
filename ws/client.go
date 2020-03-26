package ws

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
)

var (
	ErrShutdown = errors.New("websocket connection is shutdown by client")
)

func InflateDecode(in []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(in))
	defer r.Close()
	return ioutil.ReadAll(r)
}

func ZlibDecode(in []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}

func GzipDecode(in []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}

type Client struct {
	wss    string
	conn   *websocket.Conn
	closed *atomic.Bool
}

func New(ctx context.Context, wss, proxy string) (*Client, error) {
	d := websocket.DefaultDialer

	proxy = strings.TrimSpace(proxy)
	if proxy != "" {
		url, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}

		d = &websocket.Dialer{
			Proxy: http.ProxyURL(url),
		}
	}

	conn, _, err := d.DialContext(ctx, wss, nil)
	if err != nil {
		return nil, err
	}

	c := &Client{
		wss:    wss,
		conn:   conn,
		closed: atomic.NewBool(false),
	}

	return c, nil
}

func (c *Client) Write(bytes []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, bytes)
	if c.closed.Load() {
		return ErrShutdown
	}
	return err
}

func (c *Client) WriteObj(v interface{}) error {
	bytes, err := json.Marshal(v)
	if c.closed.Load() {
		return ErrShutdown
	}
	if err != nil {
		return err
	}
	return c.Write(bytes)
}

func (c *Client) Read() ([]byte, error) {
	_, bytes, err := c.conn.ReadMessage()
	if c.closed.Load() {
		return nil, ErrShutdown
	}
	return bytes, err
}

func (c *Client) ReadObj(v interface{}) error {
	_, bytes, err := c.conn.ReadMessage()
	if c.closed.Load() {
		return ErrShutdown
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func (c *Client) Close() error {
	c.closed.Store(true)
	return c.conn.Close()
}
