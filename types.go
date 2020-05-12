package auramq

//Topic define topic by string
type Topic string

//Message content of message
type Message struct {
	Topic   Topic  `json: "topic"`
	Content []byte `json:"content"`
}

//SubscribeMsg message of subscribing topics
type SubscribeMsg struct {
	Topics []Topic `json:"topics"`
}

//RegMsg register topic for subscriber
type RegMsg struct {
	topics []Topic
	client Subscriber
}

//NewRegMsg create a new register message
func NewRegMsg(client Subscriber, topics []Topic) *RegMsg {
	return &RegMsg{
		topics: topics,
		client: client,
	}
}
