package embed

import (
	"sync"

	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/msg"
)

//Broker embed broker
type Broker struct {
	router   *auramq.Router
	auth     bool
	authFunc func(*msg.AuthReq) bool
	lock     sync.RWMutex
}

//NewBroker create new embeded broker
func NewBroker(router *auramq.Router, auth bool, authFunc func(*msg.AuthReq) bool) auramq.Broker {
	return &Broker{
		router:   router,
		auth:     auth,
		authFunc: authFunc,
	}
}

//Run start embeded broker
func (broker *Broker) Run() {

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

}
