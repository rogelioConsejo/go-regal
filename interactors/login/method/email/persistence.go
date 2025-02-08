package email

import (
	"github.com/rogelioConsejo/go-regal/entities/user"
	"github.com/rogelioConsejo/go-regal/interactors/login"
	"time"
)

type Persistence interface {
	SaveLoginToken(user user.User, token HashedToken, expiration time.Time) error
	GetLoginToken(user.Name) (token HashedToken, expiration time.Time, err error)
}

type Token login.Credential
type HashedToken login.HashedCredential
