package message

type Client interface {
	Send(Address, Message) error
}

type EmailClient Client

type Address string
type Message struct {
	Subject string
	Body    string
}
