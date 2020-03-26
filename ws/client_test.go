package ws

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type dummy struct {
	Name string "json:`name`"
	Id   int    "json:`id`"
}

func TestClientMock(t *testing.T) {
	t.Skip()

	var upgrader = websocket.Upgrader{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, msg)
			if err != nil {
				break
			}
		}
	}))

	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := strings.Replace(s.URL, "http", "ws", 1)

	ctx, cancel := context.WithCancel(context.Background())
	c, err := New(ctx, u, "")
	assert.NoError(t, err)
	defer c.Close()
	defer cancel()

	msg := "test string"
	assert.NoError(t, c.Write([]byte(msg)))
	bytes, err := c.Read()
	assert.NoError(t, err)
	assert.Equal(t, msg, string(bytes))

	sent := &dummy{Name: "test_name", Id: 999}
	assert.NoError(t, c.WriteObj(sent))
	var recv dummy
	assert.NoError(t, c.ReadObj(&recv))
	assert.Equal(t, sent.Name, recv.Name)
	assert.Equal(t, sent.Id, recv.Id)
}

func TestClient(t *testing.T) {
	t.Skip()

	c, err := New(context.Background(), "wss://ws.zt.com/ws", "")
	assert.NoError(t, err)

	assert.NoError(t, c.Write([]byte("{\"method\":\"server.auth\",\"params\":[\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjExNTkyMDMsIkxvZ2luVmVyaWZ5IjoxLCJleHAiOjE1NjgyMDcwMTV9.pigLSmGlZ0IgZTRLQjc8xDaRILBduqK-dj45yRd1_iw\",\"/\"],\"id\":8417}")))
	bytes, err := c.Read()
	assert.NoError(t, err)
	fmt.Printf("resp %v\n", string(bytes))

	//assert.NoError(t, c.Write([]byte("{\"method\": \"depth.subscribe\", \"params\": [\"BTC_CNT\"], \"id\": 8418}")))
	assert.NoError(t, c.Write([]byte("{\"method\":\"depth.subscribe\",\"params\":[\"SIPC_CNT\",50,\"0.00000001\"],\"id\":5670}")))
	for true {
		bytes, err = c.Read()
		assert.NoError(t, err)
		fmt.Printf("resp %v\n", string(bytes))
	}
}
