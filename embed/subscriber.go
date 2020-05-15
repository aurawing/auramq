package embed

import (
	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/msg"
)

//Subscriber for embed broker
type Subscriber struct {
	id        string
	router    *auramq.Router
	receiver  chan *msg.Message
	Sender    chan *msg.Message
	Publisher chan *msg.Message
	rdone     chan struct{}
	wdone     chan struct{}
}

//NewSubscriber create a new embed subscriber
func NewSubscriber(id string, router *auramq.Router, receiverBufferSize, senderBufferSize, publisherBufferSize int) auramq.Subscriber {
	return &Subscriber{
		id:        id,
		router:    router,
		receiver:  make(chan *msg.Message, receiverBufferSize),
		Sender:    make(chan *msg.Message, senderBufferSize),
		Publisher: make(chan *msg.Message, publisherBufferSize),
		rdone:     make(chan struct{}),
		wdone:     make(chan struct{}),
	}
}

//ID of subscriber
func (s *Subscriber) ID() string {
	return s.id
}

//Send send message
func (s *Subscriber) Send(msg *msg.Message) bool {
	select {
	case s.receiver <- msg:
		return true
	default:
		return false
	}
}

//Run start subscriber
func (s *Subscriber) Run() {
	go func() {
	OUTER1:
		for {
			select {
			case <-s.wdone:
				break OUTER1
			case msg := <-s.receiver:
				s.Publisher <- msg
			}
		}
	}()
	go func() {
	OUTER2:
		for {
			select {
			case <-s.rdone:
				break OUTER2
			case msg := <-s.Sender:
				msg.Sender = s.ID()
				s.router.Publish(msg)
			}
		}
	}()
}

//Close close subscriber
func (s *Subscriber) Close() {
	s.router.UnregisterSubscriber(s)
	s.rdone <- struct{}{}
	s.wdone <- struct{}{}
	close(s.rdone)
	close(s.wdone)
	close(s.receiver)
	close(s.Sender)
	close(s.Publisher)
}
