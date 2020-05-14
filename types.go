package auramq

//message type
const (
	_         = iota
	P2P       //P2P message type
	BROADCAST //broadcast message type
)

//RegMsg register topic for subscriber
type RegMsg struct {
	topics []string
	client Subscriber
}

//NewRegMsg create a new register message
func NewRegMsg(client Subscriber, topics []string) *RegMsg {
	return &RegMsg{
		topics: topics,
		client: client,
	}
}
