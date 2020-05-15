package embed

import (
	"log"

	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/msg"
)

//HandShakeMsg handshake with embed broker
type HandShakeMsg struct {
	AuthMsg      *msg.AuthReq
	SubscribeMsg *msg.SubscribeReq
	HookCh       chan auramq.Subscriber
}

//Broker embed broker
type Broker struct {
	router              *auramq.Router
	auth                bool
	authFunc            func(*msg.AuthReq) bool
	HandshakeCh         chan *HandShakeMsg
	receiverBufferSize  int
	senderBufferSize    int
	publisherBufferSize int
}

//NewBroker create new embeded broker
func NewBroker(router *auramq.Router, auth bool, authFunc func(*msg.AuthReq) bool, receiverBufferSize, senderBufferSize, publisherBufferSize int) auramq.Broker {
	if receiverBufferSize == 0 {
		receiverBufferSize = 1024
	}
	if senderBufferSize == 0 {
		senderBufferSize = 1024
	}
	if publisherBufferSize == 0 {
		publisherBufferSize = 1024
	}
	return &Broker{
		router:              router,
		auth:                auth,
		authFunc:            authFunc,
		HandshakeCh:         make(chan *HandShakeMsg),
		receiverBufferSize:  receiverBufferSize,
		senderBufferSize:    senderBufferSize,
		publisherBufferSize: publisherBufferSize,
	}
}

//Run start embeded broker
func (broker *Broker) Run() {
	go func() {
	OUTER:
		for {
		INNER:
			select {
			case hs, ok := <-broker.HandshakeCh:
				if !ok {
					break OUTER
				}
				log.Println("received handshake message")
				if hs.AuthMsg.Id == "" {
					log.Println("client ID is empty")
					close(hs.HookCh)
					break INNER
				}
				if broker.auth && !broker.authFunc(hs.AuthMsg) {
					log.Println("auth failed")
					close(hs.HookCh)
					break INNER
				}
				subscriber := NewSubscriber(hs.AuthMsg.Id, broker.router, broker.receiverBufferSize, broker.senderBufferSize, broker.publisherBufferSize)
				err := broker.router.Register(subscriber, hs.SubscribeMsg.Topics)
				if err != nil {
					log.Println("client ID conflicted")
					close(hs.HookCh)
					break INNER
				}
				hs.HookCh <- subscriber
				subscriber.Run()
				log.Println("subscriber created:", subscriber.ID())
			}

		}
	}()
}

//NeedAuth if need auth when subscribe
func (broker *Broker) NeedAuth() bool {
	return broker.auth
}

//Auth authencate when subscribing
func (broker *Broker) Auth(authMsg *msg.AuthReq) bool {
	return broker.authFunc(authMsg)
}

//Close close broker
func (broker *Broker) Close() {
	close(broker.HandshakeCh)
}
