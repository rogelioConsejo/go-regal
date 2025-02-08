package login

import (
	"errors"
	"fmt"
	"github.com/rogelioConsejo/go-regal/entities/user"
	"github.com/rogelioConsejo/go-regal/interactors/message"
	"github.com/rogelioConsejo/golibs/helpers"
)

func NewUserRegistry(p RegistryPersistence, cl message.Client) UserRegistry {
	return userRegistry{
		persistence: p,
		client:      cl,
	}
}

// UserRegistry is an interface that defines the methods that a user registry should implement.
type UserRegistry interface {
	CreateUser(u user.User) error
	UserExists(u user.Name) (bool, error)
	ConfirmUserEmail(name user.Name, code ConfirmationCode) error
	UserEmailIsConfirmed(name user.Name) (bool, error)
	GetUserEmail(name user.Name) (user.Email, error)
}

type ConfirmationCode string

type userRegistry struct {
	defaultAuthenticationMethod AuthenticationMethod
	persistence                 RegistryPersistence
	client                      message.Client
}

func (ur userRegistry) GetUserEmail(name user.Name) (user.Email, error) {
	confirmed, err := ur.persistence.IsEmailConfirmed(name)
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", errors.Join(ErrGettingUserEmail, ErrEmailNotConfirmed)
	}

	email, err := ur.persistence.GetUserEmail(name)
	if err != nil {
		return "", errors.Join(ErrGettingUserEmail, err)
	}
	return email, nil
}

func (ur userRegistry) ConfirmUserEmail(name user.Name, providedCode ConfirmationCode) error {
	exists, err := ur.UserExists(name)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return errors.New("user does not exist")
	}

	savedCode, err := ur.persistence.GetConfirmationCode(name)
	if err != nil {
		return fmt.Errorf("error getting confirmation providedCode: %w", err)
	}

	if savedCode != providedCode {
		return errors.New("invalid confirmation providedCode")
	}

	err = ur.persistence.MarkEmailAsConfirmed(name)
	if err != nil {
		return fmt.Errorf("error marking email as confirmed: %w", err)
	}

	return nil
}

func (ur userRegistry) UserEmailIsConfirmed(name user.Name) (bool, error) {
	exists, err := ur.UserExists(name)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return false, errors.New("user does not exist")
	}

	return ur.persistence.IsEmailConfirmed(name)
}

func (ur userRegistry) CreateUser(u user.User) error {
	exists, err := ur.UserExists(u.Name())
	if err != nil {
		return errors.Join(ErrOnUserCreation, err)
	}
	if exists {
		return errors.Join(ErrOnUserCreation, ErrUserAlreadyExists)
	}
	err = ur.persistence.SaveUser(u)
	if err != nil {
		return errors.Join(ErrOnUserCreation, ErrCouldNotSaveUser, err)
	}
	code := generateConfirmationCode()
	err = ur.persistence.SaveConfirmationCode(u.Name(), code)
	if err != nil {
		return errors.Join(ErrOnUserCreation, ErrSavingConfirmationCode, err)
	}
	m, err := ur.getConfirmationEmail(u.Name(), code)
	if err != nil {
		return errors.Join(ErrOnUserCreation, err)
	}
	err = ur.client.Send(message.Address(u.Email()), m)
	return nil
}

func (ur userRegistry) getConfirmationEmail(name user.Name, code ConfirmationCode) (message.Message, error) {
	link, err := ur.makeConfirmationLink(name, code)
	if err != nil {
		return message.Message{}, errors.Join(ErrGeneratingConfirmationEmail, err)
	}
	m := message.Message{
		Subject: fmt.Sprintf("Hi %s, welcome to go-regal, please confirm your email", name),
		Body: fmt.Sprintf("Hi %s, welcome to go-regal, please confirm your email by clicking on the following "+
			"link: %s", name, link),
	}
	return m, nil
}

func (ur userRegistry) UserExists(u user.Name) (bool, error) {
	exists, err := ur.persistence.UserWasSaved(u)
	if err != nil {
		return false, errors.Join(ErrCouldNotCheckUser, err)
	}
	return exists, nil
}

func (ur userRegistry) makeConfirmationLink(name user.Name, confirmationCode ConfirmationCode) (string, error) {
	return fmt.Sprintf("http://localhost:8080/confirm/%s/%s", name, confirmationCode), nil
}

func generateConfirmationCode() ConfirmationCode {
	return ConfirmationCode(helpers.MakeRandomString(10))
}

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrCouldNotSaveUser = errors.New("could not save user")
var ErrCouldNotCheckUser = errors.New("could not check user")
var ErrOnUserCreation = errors.New("error on user creation")
var ErrSavingConfirmationCode = errors.New("error saving confirmation code")
var ErrGeneratingConfirmationEmail = errors.New("error generating confirmation email")
var ErrGettingUserEmail = errors.New("error getting user email")
var ErrEmailNotConfirmed = errors.New("email not confirmed")
