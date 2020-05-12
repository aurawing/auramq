package client

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

//Client websocket client
type Client struct {
	conn *websocket.Conn
}

func NewClient(wsurl string) (*Client, error) {
	u, err := url.Parse(wsurl)
	if err != nil {
		log.Printf("parse URL %s failed: %s\n", wsurl, err)
		return nil, err
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("connect URL %s failed: %s\n", wsurl, err)
		return nil, err
	}
}
