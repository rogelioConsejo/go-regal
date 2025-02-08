package login

import (
	"crypto/sha256"
	"fmt"
)

func HashCredential(cred Credential) HashedCredential {
	hasher := sha256.New()
	hasher.Write([]byte(cred))
	return HashedCredential(fmt.Sprintf("%x", hasher.Sum(nil)))
}

// Credential is a type that represents the credential that the user should provide to authenticate. Something to be
// validated against the one associated with the user entity.
type Credential string
type HashedCredential Credential

func (h HashedCredential) Check(c Credential) bool {
	return h == HashCredential(c)
}
