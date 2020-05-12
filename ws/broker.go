package ws

import (
	"log"
	"net/http"

	"github.com/aurawing/auramq"
	"github.com/gorilla/websocket"
)

//WsBroker websocket broker
type WsBroker struct {
	server               *http.Server
	router               *auramq.Router
	addr                 string
	auth                 bool
	authFunc             func([]byte) bool
	readBufferSize       int
	writeBufferSize      int
	subscriberBufferSize int
	pongWait             int
	writeWait            int
}

//NewWsBroker create new websocket broker
func NewWsBroker(router *auramq.Router, addr string, auth bool, authFunc func([]byte) bool, subscriberBufferSize, readBufferSize, writeBufferSize, pongWait, writeWait int) auramq.Broker {
	if subscriberBufferSize == 0 {
		subscriberBufferSize = 64
	}
	if readBufferSize == 0 {
		readBufferSize = 4096
	}
	if writeBufferSize == 0 {
		writeBufferSize = 4096
	}
	if pongWait == 0 {
		pongWait = 60
	}
	if writeWait == 0 {
		writeWait = 10
	}
	return &WsBroker{
		router:               router,
		addr:                 addr,
		auth:                 auth,
		authFunc:             authFunc,
		readBufferSize:       readBufferSize,
		writeBufferSize:      writeBufferSize,
		subscriberBufferSize: subscriberBufferSize,
		pongWait:             pongWait,
		writeWait:            writeWait,
	}
}

//Run start websocket broker
func (broker *WsBroker) Run() {
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
		if broker.NeedAuth() {
			_, b, err := conn.ReadMessage()
			if err != nil {
				log.Println("read auth message failed:", err)
				// conn.WriteMessage(websocket.CloseMessage, []byte{})
				conn.Close()
				return
			}
			// if msgType != websocket.BinaryMessage {
			// 	log.Println("auth message should be binary format")
			// 	return
			// }
			if !broker.Auth(b) {
				log.Println("auth failed")
				// conn.WriteMessage(websocket.CloseMessage, []byte{})
				conn.Close()
				return
			}
		}
		subscribeMsg := new(auramq.SubscribeMsg)
		if err := conn.ReadJSON(subscribeMsg); err != nil {
			log.Printf("error when decode topics for subscribing: %s\n", err)
			// conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()
			return
		}
		subscriber := NewWsSubscriber(broker.router, conn, broker.subscriberBufferSize, broker.pongWait, broker.writeWait)
		subscriber.Run()
		broker.router.Register(subscriber, subscribeMsg.Topics)
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("httpserver: ListenAndServe() error: %s", err)
		}
	}()
	broker.server = srv
}

//NeedAuth if need auth when subscribe
func (broker *WsBroker) NeedAuth() bool {
	return broker.auth
}

//Auth authencate when subscribing
func (broker *WsBroker) Auth(authMsg []byte) bool {
	return broker.authFunc(authMsg)
}

//Close close http server
func (broker *WsBroker) Close() {
	if err := broker.server.Shutdown(nil); err != nil {
		log.Printf("httpserver: Shutdown() error: %s", err)
	}
}
