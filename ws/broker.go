package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/msg"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

//Broker websocket broker
type Broker struct {
	server               *http.Server
	router               *auramq.Router
	addr                 string
	auth                 bool
	authFunc             func(*msg.AuthReq) bool
	readBufferSize       int
	writeBufferSize      int
	subscriberBufferSize int
	pingWait             int
	readWait             int
	writeWait            int
}

//NewBroker create new websocket broker
func NewBroker(router *auramq.Router, addr string, auth bool, authFunc func(*msg.AuthReq) bool, subscriberBufferSize, readBufferSize, writeBufferSize, pingWait, readWait, writeWait int) auramq.Broker {
	if subscriberBufferSize == 0 {
		subscriberBufferSize = 1024
	}
	if readBufferSize == 0 {
		readBufferSize = 4096
	}
	if writeBufferSize == 0 {
		writeBufferSize = 4096
	}
	if pingWait == 0 {
		pingWait = 30
	}
	if readWait == 0 {
		readWait = 60
	}
	if writeWait == 0 {
		writeWait = 10
	}
	return &Broker{
		router:               router,
		addr:                 addr,
		auth:                 auth,
		authFunc:             authFunc,
		readBufferSize:       readBufferSize,
		writeBufferSize:      writeBufferSize,
		subscriberBufferSize: subscriberBufferSize,
		pingWait:             pingWait,
		readWait:             readWait,
		writeWait:            writeWait,
	}
}

//Run start websocket broker
func (broker *Broker) Run() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:    broker.readBufferSize,
		WriteBufferSize:   broker.writeBufferSize,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	srv := &http.Server{Addr: broker.addr}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("subscribe error:", err)
			return
		}
		//receive auth message
		conn.SetReadDeadline(time.Now().Add(time.Duration(broker.readWait) * time.Second))
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Println("read auth request failed:", err)
			conn.Close()
			return
		}
		authReq := new(msg.AuthReq)
		err = proto.Unmarshal(b, authReq)
		if err != nil {
			log.Println("unmarshal auth request failed:", err)
			conn.Close()
			return
		}
		conn.SetWriteDeadline(time.Now().Add(time.Duration(broker.writeWait) * time.Second))
		if authReq.Id == "" {
			log.Printf("client ID is empty")
			authAck := &msg.Ack{Ack: false}
			b, err := proto.Marshal(authAck)
			if err != nil {
				log.Println("unmarshal auth ack failed:", err)
				conn.Close()
				return
			}
			err = conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				log.Println("write auth ack failed:", err)
				conn.Close()
				return
			}
			log.Println("auth failed")
			conn.Close()
			return
		}
		if broker.NeedAuth() && !broker.Auth(authReq) {
			log.Printf("auth failed")
			authAck := &msg.Ack{Ack: false}
			b, err := proto.Marshal(authAck)
			if err != nil {
				log.Println("unmarshal auth ack failed:", err)
				conn.Close()
				return
			}
			err = conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				log.Println("write auth ack failed:", err)
				conn.Close()
				return
			}
			log.Println("auth failed")
			conn.Close()
			return
		}
		authAck := &msg.Ack{Ack: true}
		b, err = proto.Marshal(authAck)
		if err != nil {
			log.Println("unmarshal auth ack failed:", err)
			conn.Close()
			return
		}
		err = conn.WriteMessage(websocket.BinaryMessage, b)
		if err != nil {
			log.Println("write auth ack failed:", err)
			conn.Close()
			return
		}
		//receive subscribe message
		conn.SetReadDeadline(time.Now().Add(time.Duration(broker.readWait) * time.Second))
		_, b, err = conn.ReadMessage()
		if err != nil {
			log.Println("error when read topics for subscribing:", err)
			conn.Close()
			return
		}
		subscribeReq := new(msg.SubscribeReq)
		err = proto.Unmarshal(b, subscribeReq)
		if err != nil {
			log.Println("unmarshal topics for subscribing failed:", err)
			conn.Close()
			return
		}

		subscriber := NewSubscriber(authReq.Id, broker.router, conn, broker.subscriberBufferSize, broker.pingWait, broker.readWait, broker.writeWait)
		err = broker.router.Register(subscriber, subscribeReq.Topics)
		conn.SetWriteDeadline(time.Now().Add(time.Duration(broker.writeWait) * time.Second))
		if err != nil {
			log.Println("client ID conflicted")
			subAck := &msg.Ack{Ack: false}
			b, err := proto.Marshal(subAck)
			if err != nil {
				log.Println("unmarshal subscribe ack failed:", err)
				conn.Close()
				return
			}
			err = conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				log.Println("write subscribe ack failed:", err)
				conn.Close()
				return
			}
			conn.Close()
			return
		}
		subAck := &msg.Ack{Ack: true}
		b, err = proto.Marshal(subAck)
		if err != nil {
			log.Println("unmarshal subscribe ack failed:", err)
			conn.Close()
			return
		}
		err = conn.WriteMessage(websocket.BinaryMessage, b)
		if err != nil {
			log.Println("write subscribe ack failed:", err)
			conn.Close()
			return
		}
		subscriber.Run()
		log.Println("subscriber created:", subscriber.ID())
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("httpserver: ListenAndServe() error: %s", err)
		}
	}()
	broker.server = srv
}

//NeedAuth if need auth when subscribe
func (broker *Broker) NeedAuth() bool {
	return broker.auth
}

//Auth authencate when subscribing
func (broker *Broker) Auth(authMsg *msg.AuthReq) bool {
	return broker.authFunc(authMsg)
}

//Close close http server
func (broker *Broker) Close() {
	if err := broker.server.Shutdown(nil); err != nil {
		log.Printf("httpserver: Shutdown() error: %s", err)
	}
	//broker.router.Close()
}
