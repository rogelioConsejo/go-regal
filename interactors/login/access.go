package login

import (
	"errors"
	"github.com/rogelioConsejo/go-regal/entities/user"
	"github.com/rogelioConsejo/golibs/helpers"
)

func NewAccess(dam AuthenticationMethod) Access {
	return access{
		defaultAuthenticationMethod: dam,
	}
}

type AccessToken Credential

type Access interface {
	PrepareAuthentication(user.User) error
	PerformAuthentication(user.User, Credential) (AccessToken, error)
}

type access struct {
	defaultAuthenticationMethod AuthenticationMethod
}

func (a access) PrepareAuthentication(u user.User) error {
	err := a.defaultAuthenticationMethod.SetupAuthenticationAttempt(u)
	if err != nil {
		return errors.Join(ErrPreparingAuthenticationAttempt, err)
	}
	return nil
}

func (a access) PerformAuthentication(u user.User, credential Credential) (AccessToken, error) {
	isValid, err := a.defaultAuthenticationMethod.Authenticate(u, credential)
	if err != nil {
		return "", errors.Join(ErrPerformingAuthentication, err)
	}
	if !isValid {
		return "", nil
	}
	atk := generateKey()
	return atk, nil
}

func generateKey() AccessToken {
	return AccessToken(helpers.MakeRandomString(10))
}

var ErrPreparingAuthenticationAttempt = errors.New("error preparing authentication attempt")
var ErrPerformingAuthentication = errors.New("error performing authentication")
