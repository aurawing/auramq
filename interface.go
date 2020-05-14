package auramq

import "github.com/aurawing/auramq/msg"

//Broker interface
type Broker interface {
	Run()
	Close()
	NeedAuth() bool
	Auth(auth *msg.AuthReq) bool
}

//Subscriber interface
type Subscriber interface {
	ID() string
	Send(msg *msg.Message) bool
	Run()
	Close()
}

//Client interface
type Client interface {
	Send(to string, content []byte) bool
	Publish(topic string, content []byte) bool
	Close()
}
