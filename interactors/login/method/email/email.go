package email

import (
	"errors"
	"github.com/rogelioConsejo/go-regal/entities/user"
	"github.com/rogelioConsejo/go-regal/interactors/login"
	"github.com/rogelioConsejo/golibs/helpers"
	"time"
)

func GetEmailMethod(p Persistence, s Sender) login.AuthenticationMethod {
	return &emailMethod{
		expiration:  15 * time.Minute,
		emailSender: s,
		persistence: p,
	}
}

type emailMethod struct {
	persistence Persistence
	emailSender Sender
	expiration  time.Duration
}

func (e emailMethod) SetupAuthenticationAttempt(u user.User) error {
	tk := Token(helpers.MakeRandomString(10))
	hashedToken := login.HashCredential(login.Credential(tk))
	err := e.persistence.SaveLoginToken(u, HashedToken(hashedToken), time.Now().Add(e.expiration))
	if err != nil {
		return errors.Join(ErrSavingLoginToken, err)
	}
	err = e.emailSender.sendToken(u.Email(), tk)
	if err != nil {
		return errors.Join(ErrSendingLoginToken, err)
	}
	return nil
}

func (e emailMethod) Authenticate(u user.User, credential login.Credential) (bool, error) {
	hashedToken, tokenExpiration, err := e.persistence.GetLoginToken(u.Name())
	if err != nil {
		return false, errors.Join(ErrAuthenticating, ErrRetrievingLoginToken, err)
	}
	if time.Now().After(tokenExpiration) {
		return false, nil
	}

	return login.HashedCredential(hashedToken).Check(credential), nil
}

var ErrSavingLoginToken = errors.New("error saving login token")
var ErrSendingLoginToken = errors.New("error sending login token")
var ErrAuthenticating = errors.New("could not authenticate user")
var ErrRetrievingLoginToken = errors.New("failed to retrieve login token")
