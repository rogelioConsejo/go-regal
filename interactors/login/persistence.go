package login

import "github.com/rogelioConsejo/go-regal/entities/user"

type RegistryPersistence interface {
	SaveUser(u user.User) error
	UserWasSaved(u user.Name) (bool, error)
	SaveConfirmationCode(u user.Name, c ConfirmationCode) error
}

type AccessPersistence interface {
}

type DoorLockPersistence interface {
}
