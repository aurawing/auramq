package auramq

//Broker interface
type Broker interface {
	Run()
	Close()
	NeedAuth() bool
	Auth(auth []byte) bool
}

//Subscriber interface
type Subscriber interface {
	Send(msg *Message) bool
	Run()
	Close()
}
