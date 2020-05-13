package client

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/aurawing/auramq/msg"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

//Client websocket client
type Client struct {
	URL          *url.URL
	conn         *websocket.Conn
	CallbackFunc func(*msg.Message)
	authMsg      *msg.AuthReq
	receiver     chan *msg.Message
	pingWait     int
	readWait     int
	writeWait    int
}

//Connect create a new websocket client and connect to server
func Connect(wsurl string, callback func(*msg.Message), authMsg *msg.AuthReq, topics []string, receiverBufferSize, pingWait, readWait, writeWait int) (*Client, error) {
	if receiverBufferSize == 0 {
		receiverBufferSize = 64
	}
	if pingWait == 0 {
		pingWait = 30
	}
	if readWait == 0 {
		readWait = 60
	}
	if writeWait == 0 {
		writeWait = 10
	}
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
	if authMsg != nil {
		b, err := proto.Marshal(authMsg)
		if err != nil {
			log.Println("unmarshal auth failed:", err)
			conn.Close()
			return nil, err
		}
		conn.SetWriteDeadline(time.Now().Add(time.Duration(writeWait) * time.Second))
		err = conn.WriteMessage(websocket.BinaryMessage, b)
		if err != nil {
			log.Println("client auth failed:", err)
			conn.Close()
			return nil, err
		}
		conn.SetReadDeadline(time.Now().Add(time.Duration(readWait) * time.Second))
		_, b, err = conn.ReadMessage()
		if err != nil {
			log.Println("error when read topics for subscribing:", err)
			conn.Close()
			return nil, err
		}
		authAck := new(msg.Ack)
		err = proto.Unmarshal(b, authAck)
		if err != nil {
			log.Println("unmarshal topics for subscribing failed:", err)
			conn.Close()
			return nil, err
		}
		if authAck.Ack {
			log.Println("auth success")
		} else {
			log.Println("auth failed")
			conn.Close()
			return nil, errors.New("auth failed")
		}
	}
	b, err := proto.Marshal(&msg.SubscribeReq{Topics: topics})
	if err != nil {
		log.Println("marshal subscribe request failed")
		conn.Close()
		return nil, err
	}
	conn.SetWriteDeadline(time.Now().Add(time.Duration(writeWait) * time.Second))
	err = conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		log.Println("send subscribe request failed")
		conn.Close()
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(time.Duration(readWait) * time.Second))
	_, b, err = conn.ReadMessage()
	if err != nil {
		log.Println("read subscribe ack failed")
		conn.Close()
		return nil, err
	}
	subAck := new(msg.Ack)
	err = proto.Unmarshal(b, subAck)
	if err != nil {
		log.Println("unmarshal subscribe ack failed:", err)
		conn.Close()
		return nil, err
	}
	if subAck.Ack {
		log.Println("subscribe topics success")
	} else {
		log.Println("subscribe topics failed")
		conn.Close()
		return nil, errors.New("subscribe topics failed")
	}
	client := &Client{URL: u, conn: conn, CallbackFunc: callback, receiver: make(chan *msg.Message, receiverBufferSize), authMsg: authMsg, pingWait: pingWait, readWait: readWait, writeWait: writeWait}
	return client, nil
}

//Run client
func (c *Client) Run() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.readWait) * time.Second))
		c.conn.SetPongHandler(func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.readWait) * time.Second))
			return nil
		})
		for {
			_, b, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("unexpected close error: %s", err)
				} else {
					log.Printf("read message error: %s", err)
				}
				break
			}
			msg := new(msg.Message)
			err = proto.Unmarshal(b, msg)
			if err != nil {
				fmt.Println("unmarshal message failed")
				continue
			}
			c.CallbackFunc(msg)
		}
		c.Close()
		wg.Done()
	}()

	go func() {
		ticker := time.NewTicker(time.Duration(c.pingWait) * time.Second)
		defer ticker.Stop()
	OUT:
		for {
			select {
			case <-ticker.C:
				c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.writeWait) * time.Second))
				err := c.conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Println("write ping message error:", err)
					break
				}
			case msg, ok := <-c.receiver:
				if !ok {
					log.Println("reciver channel closed")
					break OUT
				}
				b, err := proto.Marshal(msg)
				if err != nil {
					log.Println("marshal message error:", err)
					break OUT
				}
				c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.writeWait) * time.Second))
				err = c.conn.WriteMessage(websocket.BinaryMessage, b)
				if err != nil {
					log.Println("write message error:", err)
					break OUT
				}
			}
		}
		c.conn.Close()
		wg.Done()
	}()
	wg.Wait()
	log.Println("connection closed")
}

//Publish one message
func (c *Client) Publish(message *msg.Message) bool {
	select {
	case c.receiver <- message:
		return true
	default:
		return false
	}
}

//Close client
func (c *Client) Close() {
	close(c.receiver)
	c.conn.Close()
}
