package ws

import (
	"log"
	"time"

	"github.com/aurawing/auramq"
	"github.com/gorilla/websocket"
)

//WsSubscriber subscriber for websocket broker
type WsSubscriber struct {
	router    *auramq.Router
	conn      *websocket.Conn
	receiver  chan *auramq.Message
	pongWait  int
	writeWait int
}

//NewWsSubscriber create a new websocket subscriber
func NewWsSubscriber(router *auramq.Router, conn *websocket.Conn, subscriberBufferSize, pongWait, writeWait int) auramq.Subscriber {
	return &WsSubscriber{
		router:    router,
		conn:      conn,
		receiver:  make(chan *auramq.Message, subscriberBufferSize),
		pongWait:  pongWait,
		writeWait: writeWait,
	}
}

//Send send message
func (s *WsSubscriber) Send(msg *auramq.Message) bool {
	select {
	case s.receiver <- msg:
		return true
	default:
		return false
	}
}

//Run start subscriber
func (s *WsSubscriber) Run() {
	go s.readPump()
	go s.writePump()
}

//Close close subscriber
func (s *WsSubscriber) Close() {
	s.conn.Close()
}

func (s *WsSubscriber) readPump() {
	defer func() {
		s.router.UnregisterSubscriber(s)
		s.Close()
		close(s.receiver)
	}()
	//s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(time.Duration(s.pongWait) * time.Second))
	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(time.Duration(s.pongWait) * time.Second))
		return nil
	})
	for {
		//_, message, err := s.conn.ReadMessage()
		publishMsg := new(auramq.Message)
		err := s.conn.ReadJSON(publishMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			} else {
				log.Printf("read message error: %v", err)
			}
			break
		}
		s.router.Publish(publishMsg)
	}
}

func (s *WsSubscriber) writePump() {
	for {
		select {
		case message, ok := <-s.receiver:
			if !ok {
				return
			}
			s.conn.SetWriteDeadline(time.Now().Add(time.Duration(s.writeWait) * time.Second))
			err := s.conn.WriteJSON(message)
			if err != nil {
				s.conn.Close()
				return
			}
		}
	}
}
