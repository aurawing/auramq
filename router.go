package auramq

import (
	"sync"

	"github.com/fatih/set"
)

//Router routing message to subscriber
type Router struct {
	rtable    map[Topic]set.Interface
	rrtable   map[Subscriber]set.Interface
	broadcast chan *Message
	done      chan struct{}
	lock      sync.RWMutex
}

//NewRouter create a new router instance
func NewRouter(bufferSize int) *Router {
	return &Router{
		rtable:    make(map[Topic]set.Interface),
		rrtable:   make(map[Subscriber]set.Interface),
		broadcast: make(chan *Message, bufferSize),
		done:      make(chan struct{}),
	}
}

//Register register topics for subscriber
func (router *Router) Register(client Subscriber, topics []Topic) {
	router.lock.Lock()
	defer router.lock.Unlock()
	if router.rrtable[client] == nil {
		router.rrtable[client] = set.New(set.NonThreadSafe)
	}
	s := set.New(set.NonThreadSafe)
	for _, t := range topics {
		s.Add(t)
	}
	intersect := set.Intersection(s, router.rrtable[client])
	for _, t := range intersect.List() {
		router.rrtable[client].Add(t)
	}
	for _, topic := range intersect.List() {
		if router.rtable[topic.(Topic)] == nil {
			router.rtable[topic.(Topic)] = set.New(set.NonThreadSafe)
		}
		router.rtable[topic.(Topic)].Add(client)
	}
}

//UnregisterSubscriber unregister all topics for subscriber
func (router *Router) UnregisterSubscriber(client Subscriber) {
	router.lock.Lock()
	defer router.lock.Unlock()
	if _, ok := router.rrtable[client]; !ok {
		return
	}
	if router.rrtable[client].Size() == 0 {
		delete(router.rrtable, client)
		return
	}
	topics := router.rrtable[client].List()
	topicList := make([]Topic, 0)
	for _, t := range topics {
		topicList = append(topicList, t.(Topic))
	}
	router.unregister(client, topicList)
}

//Unregister unregister topics for subscriber
func (router *Router) Unregister(client Subscriber, topics []Topic) {
	router.lock.Lock()
	defer router.lock.Unlock()
	router.unregister(client, topics)
}

func (router *Router) unregister(client Subscriber, topics []Topic) {
	if _, ok := router.rrtable[client]; !ok {
		return
	}
	if router.rrtable[client].Size() == 0 {
		delete(router.rrtable, client)
		return
	}
	for _, topic := range topics {
		router.rtable[topic].Remove(client)
		if router.rtable[topic].Size() == 0 {
			delete(router.rtable, topic)
		}
		router.rrtable[client].Remove(topic)
	}
	if router.rrtable[client].Size() == 0 {
		delete(router.rrtable, client)
	}
}

//Publish publish message to a topic
func (router *Router) Publish(msg *Message) {
	router.broadcast <- msg
}

//Run start router
func (router *Router) Run() {
	for {
		select {
		case msg := <-router.broadcast:
			if router.rtable[msg.Topic] != nil {
				for _, client := range router.rtable[msg.Topic].List() {
					cli := client.(Subscriber)
					cli.Send(msg)
					// if !cli.Send(msg) {
					// 	cli.Close()
					// 	router.UnregisterSubscriber(cli)
					// }
				}
			}
		case _ = <-router.done:
			break
		}
	}
	close(router.broadcast)
	close(router.done)
	for cli := range router.rrtable {
		cli.Close()
	}
}

//Close shutdown router
func (router *Router) Close() {
	router.done <- struct{}{}
}
