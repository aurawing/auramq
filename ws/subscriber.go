package ws

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/aurawing/auramq"
	"github.com/aurawing/auramq/msg"
	"github.com/gorilla/websocket"
)

//Subscriber for websocket broker
type Subscriber struct {
	id        string
	router    *auramq.Router
	conn      *websocket.Conn
	receiver  chan *msg.Message
	pingWait  int
	readWait  int
	writeWait int
}

//NewSubscriber create a new websocket subscriber
func NewSubscriber(id string, router *auramq.Router, conn *websocket.Conn, subscriberBufferSize, pingWait, readWait, writeWait int) auramq.Subscriber {
	return &Subscriber{
		id:        id,
		router:    router,
		conn:      conn,
		receiver:  make(chan *msg.Message, subscriberBufferSize),
		pingWait:  pingWait,
		readWait:  readWait,
		writeWait: writeWait,
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
	go s.readPump()
	go s.writePump()
}

//Close close subscriber
func (s *Subscriber) Close() {
	s.conn.Close()
}

func (s *Subscriber) readPump() {
	defer func() {
		s.router.UnregisterSubscriber(s)
		s.Close()
		close(s.receiver)
	}()
	//s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(time.Duration(s.readWait) * time.Second))
	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(time.Duration(s.readWait) * time.Second))
		return nil
	})
	for {
		_, b, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %s", err)
			} else {
				log.Printf("read message error: %s", err)
			}
			break
		}
		publishMsg := new(msg.Message)
		err = proto.Unmarshal(b, publishMsg)
		if err != nil {
			log.Printf("unexpected message: %s", err)
			break
		}
		publishMsg.Sender = s.ID()
		s.router.Publish(publishMsg)
	}
}

func (s *Subscriber) writePump() {
	ticker := time.NewTicker(time.Duration(s.pingWait) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(time.Duration(s.writeWait) * time.Second))
			err := s.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("write ping message error:", err)
				break
			}
		case message, ok := <-s.receiver:
			if !ok {
				return
			}
			b, err := proto.Marshal(message)
			if err != nil {
				s.conn.Close()
				return
			}
			s.conn.SetWriteDeadline(time.Now().Add(time.Duration(s.writeWait) * time.Second))
			err = s.conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				s.conn.Close()
				return
			}
		}
	}
}
