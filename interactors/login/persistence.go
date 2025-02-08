package login

import (
	"github.com/rogelioConsejo/go-regal/entities/user"
)

type RegistryPersistence interface {
	UserPersistence
	UserConfirmationCodeRegistryPersistence
	UserEmailRegistryPersistence
}

type UserEmailRegistryPersistence interface {
	MarkEmailAsConfirmed(u user.Name) error
	IsEmailConfirmed(u user.Name) (bool, error)
	GetUserEmail(name user.Name) (user.Email, error)
}

type UserConfirmationCodeRegistryPersistence interface {
	SaveConfirmationCode(u user.Name, c ConfirmationCode) error
	GetConfirmationCode(u user.Name) (ConfirmationCode, error)
}

type UserPersistence interface {
	UserWasSaved(u user.Name) (bool, error)
	SaveUser(u user.User) error
}

type AccessPersistence interface {
}

type DoorLockPersistence interface {
}
