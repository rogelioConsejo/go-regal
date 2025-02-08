package email

import (
	"errors"
	"github.com/rogelioConsejo/go-regal/entities/user"
	"github.com/rogelioConsejo/go-regal/interactors/message"
)

func GetSender(c message.EmailClient) Sender {
	return &senderClient{c}
}

type Sender interface {
	sendToken(email user.Email, tk Token) error // may need to change to make more general
}

type senderClient struct {
	client message.Client
}

func (s senderClient) sendToken(email user.Email, tk Token) error {
	// formatting (including the message itself) will need to be handled by a separate module when we want to make
	// this configurable
	err := s.client.Send(message.Address(email), message.Message{
		Subject: "Your login token",
		Body:    "Your login token is " + string(tk),
	})

	if err != nil {
		return errors.Join(errCouldNotSendToken, err)
	}

	return nil
}

var errCouldNotSendToken = errors.New("could not send token")
