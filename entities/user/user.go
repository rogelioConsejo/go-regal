package user

import (
	"errors"
	"strings"
)

// New creates an immutable User object with a name and email.
func New(n Name, e Email) (User, error) {
	err := validateNameAndEmail(n, e)
	if err != nil {
		return nil, errors.Join(ErrCreatingUser, err)
	}

	return user{
		n,
		e,
	}, nil
}

func validateNameAndEmail(n Name, e Email) error {
	if n == "" {
		return ErrEmptyName
	}
	emailErr := validateEmail(e)
	if emailErr != nil {
		return emailErr
	}
	return nil
}

func validateEmail(e Email) error {
	if e == "" {
		return ErrEmptyEmail
	}
	if !strings.Contains(string(e), "@") {
		return errors.Join(InvalidEmail, errors.New("email must contain @ symbol"))
	}
	if strings.HasPrefix(string(e), "@") {
		return errors.Join(InvalidEmail, errors.New("email cannot start with @ symbol"))
	}
	if !strings.Contains(string(e), ".") {
		return errors.Join(InvalidEmail, errors.New("email must contain . symbol"))
	}
	if strings.LastIndex(string(e), ".") < strings.LastIndex(string(e), "@") {
		return errors.Join(InvalidEmail, errors.New("email must contain @ after the last . symbol"))
	}
	if strings.HasSuffix(string(e), ".") {
		return errors.Join(InvalidEmail, errors.New("email cannot end with . symbol"))
	}
	return nil
}

type User interface {
	Name() Name
	Email() Email
}

type user struct {
	name  Name
	email Email
}

func (u user) Email() Email {
	return u.email
}

func (u user) Name() Name {
	return u.name
}

type Name string
type Email string

var ErrCreatingUser = errors.New("error creating user")
var ErrEmptyName = errors.New("a user cannot have an empty name")
var ErrEmptyEmail = errors.New("a user cannot have an empty email")
var InvalidEmail = errors.New("email is invalid")
