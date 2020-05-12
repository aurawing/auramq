package main

import (
	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/ws"
)

func main() {
	router := auramq.NewRouter(1024)
	go router.Run()
	broker := ws.NewBroker(router, ":8080", false, nil, 0, 0, 0, 0, 0)
	broker.Run()
	ch := make(chan struct{})
	<-ch
}
