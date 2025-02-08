package login

import (
	"github.com/rogelioConsejo/go-regal/entities/user"
)

// AuthenticationMethod is an interface that defines the methods that an authentication method should implement.
//
// For example:
//
// - a password authentication method should not do anything in the SetupAuthenticationAttempt method and
// should compare the password's hash with the one stored with the user entity in the Authenticate method.
//
// - an email authentication method should send the login token to the user's email in the SetupAuthenticationAttempt
// method and should compare the token with the one stored with the user entity  in the Authenticate method (or
// its hash).
type AuthenticationMethod interface {
	// SetupAuthenticationAttempt prepares the user for the authentication attempt. It should be expected to be called
	// before the Authenticate method, even if there is nothing to do to prepare the user for their specific
	// authentication method's attempt.
	SetupAuthenticationAttempt(user.User) error
	// Authenticate validates the user's credential against the one associated with the user entity.
	Authenticate(user.User, Credential) (bool, error)
}
