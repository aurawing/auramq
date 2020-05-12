package auramq

import "github.com/aurawing/auramq/msg"

//Broker interface
type Broker interface {
	Run()
	Close()
	NeedAuth() bool
	Auth(auth []byte) bool
}

//Subscriber interface
type Subscriber interface {
	Send(msg *msg.Message) bool
	Run()
	Close()
}
