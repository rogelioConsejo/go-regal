package login

import (
	"github.com/google/uuid"
	"testing"
)

func TestHashCredential(t *testing.T) {
	var cred Credential = Credential(uuid.NewString())
	var hash HashedCredential = HashCredential(cred)
	if !hash.Check(cred) {
		t.Fatal("the hashed credential did not check")
	}
}
