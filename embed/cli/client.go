package cli

import (
	"errors"

	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/embed"
	"github.com/aurawing/auramq/msg"
)

//Client websocket client
type Client struct {
	ID           string
	CallbackFunc func(*msg.Message)
	authMsg      *msg.AuthReq
	subscriber   *embed.Subscriber
}

//Connect create a new embed client
func Connect(broker *embed.Broker, callback func(*msg.Message), authMsg *msg.AuthReq, topics []string) (*Client, error) {
	if callback == nil {
		return nil, errors.New("callback function can not be nil")
	}
	if authMsg == nil {
		return nil, errors.New("auth message can not be nil")
	}
	if authMsg.Id == "" {
		return nil, errors.New("client ID in auth message can not be empty")
	}
	if topics == nil || len(topics) == 0 {
		return nil, errors.New("topics can not be empty")
	}
	handshake := &embed.HandShakeMsg{AuthMsg: authMsg, SubscribeMsg: &msg.SubscribeReq{Topics: topics}, HookCh: make(chan auramq.Subscriber)}
	broker.HandshakeCh <- handshake
	subscriber, ok := <-handshake.HookCh
	if !ok {
		return nil, errors.New("connect failed")
	}
	close(handshake.HookCh)
	subs, ok := subscriber.(*embed.Subscriber)
	if !ok {
		return nil, errors.New("wrong type of subscriber")
	}
	client := &Client{ID: authMsg.Id, CallbackFunc: callback, authMsg: authMsg, subscriber: subs}
	return client, nil
}

//Run client
func (c *Client) Run() {
OUTER:
	for {
		select {
		case msg, ok := <-c.subscriber.Publisher:
			if !ok {
				break OUTER
			}
			c.CallbackFunc(msg)
		}
	}
}

//Send one message to another client
func (c *Client) Send(to string, content []byte) bool {
	message := &msg.Message{Type: auramq.P2P, Destination: to, Content: content}
	select {
	case c.subscriber.Sender <- message:
		return true
	default:
		return false
	}
}

//Publish one message
func (c *Client) Publish(topic string, content []byte) bool {
	message := &msg.Message{Type: auramq.BROADCAST, Destination: topic, Content: content}
	select {
	case c.subscriber.Sender <- message:
		return true
	default:
		return false
	}
}

//Close client
func (c *Client) Close() {
	c.subscriber.Close()
}
