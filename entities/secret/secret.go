package secret

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Secret string

func (s Secret) Hash() (HashedSecret, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return HashedSecret(fmt.Sprintf("%s", hashedBytes)), nil
}

type HashedSecret string

func (hs HashedSecret) Compare(s Secret) bool {
	return bcrypt.CompareHashAndPassword([]byte(hs), []byte(s)) == nil
}
