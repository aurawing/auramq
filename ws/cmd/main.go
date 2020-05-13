package main

import (
	"fmt"
	"time"

	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/msg"
	"github.com/aurawing/auramq/ws"
	client "github.com/aurawing/auramq/ws/cli"
)

func main() {
	router := auramq.NewRouter(1024)
	go router.Run()
	broker := ws.NewBroker(router, ":8080", false, nil, 0, 0, 0, 0, 0, 0)
	broker.Run()

	cli1, err := client.Connect("ws://127.0.0.1:8080/ws", callback, nil, []string{"test"}, 0, 0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		cli1.Run()
	}()

	// cli2, err := client.Connect("ws://127.0.0.1:8080/ws", callback, nil, []string{"test"}, 0, 0, 0, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// go func() {
	// 	cli2.Run()
	// }()

	b := cli1.Publish(&msg.Message{Topic: "test", Content: []byte("hahaha")})
	fmt.Println(b)
	time.Sleep(20 * time.Second)
	broker.Close()
	time.Sleep(100 * time.Second)
}

func callback(msg *msg.Message) {
	fmt.Println("received: ", string(msg.Content))
}
