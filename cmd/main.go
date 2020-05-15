package main

import (
	"fmt"
	"time"

	"github.com/aurawing/auramq/embed"

	"github.com/aurawing/auramq"
	ebclient "github.com/aurawing/auramq/embed/cli"
	"github.com/aurawing/auramq/msg"
	"github.com/aurawing/auramq/ws"
	wsclient "github.com/aurawing/auramq/ws/cli"
)

func main() {
	router := auramq.NewRouter(1024)
	go router.Run()

	wsbroker := ws.NewBroker(router, ":8080", true, auth, 0, 0, 0, 0, 0, 0)
	wsbroker.Run()

	ebbroker := embed.NewBroker(router, true, auth, 0, 0, 0)
	ebbroker.Run()

	cli1, err := wsclient.Connect("ws://127.0.0.1:8080/ws", callback, &msg.AuthReq{Id: "aaa", Credential: []byte("welcome")}, []string{"test"}, 0, 0, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		cli1.Run()
	}()

	cli2, err := ebclient.Connect(ebbroker.(*embed.Broker), callback, &msg.AuthReq{Id: "ccc", Credential: []byte("welcome")}, []string{"test"})
	go func() {
		cli2.Run()
	}()
	// cli2, err := client.Connect("ws://127.0.0.1:8080/ws", callback, &msg.AuthReq{Id: "bbb", Credential: []byte("welcome")}, []string{"test"}, 0, 0, 0, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// go func() {
	// 	cli2.Run()
	// }()

	b := cli1.Send("ccc", []byte("hahaha"))
	fmt.Println(b)
	b = cli2.Publish("test", []byte("yohoho"))
	time.Sleep(2000 * time.Second)
	wsbroker.Close()
	ebbroker.Close()
	time.Sleep(100 * time.Second)
}

func callback(msg *msg.Message) {
	fmt.Printf("from %s to %s: %s\n", msg.Sender, msg.Destination, string(msg.Content))
}

func auth(b *msg.AuthReq) bool {
	if string(b.Credential) == "welcome" {
		return true
	}
	return false
}
