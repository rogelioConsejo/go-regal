package secret

import "testing"

func TestSecret_Hash(t *testing.T) {
	var s Secret = "some secret"
	var hs HashedSecret
	hs, err := s.Hash()
	if err != nil {
		t.Fatalf("Secret.Hash should not return an error, but got %v", err)
	}
	if hs == "" {
		t.Errorf("HashedSecret should not be empty")
	}
	if !hs.Compare(s) {
		t.Errorf("HashedSecret.Compare should return true")
	}
}

func TestHashedSecret_Compare(t *testing.T) {
	var s Secret = "some secret"
	var hs HashedSecret
	hs, err := s.Hash()
	if err != nil {
		t.Fatalf("Secret.Hash should not return an error, but got %v", err)
	}
	if !hs.Compare(s) {
		t.Fatal("HashedSecret.Compare should return true")
	}
	var anotherSecret Secret = "another secret"
	if hs.Compare(anotherSecret) {
		t.Fatal("HashedSecret.Compare should return false")
	}
}

type Key Secret

func (k Key) Hash() (HashedKey, error) {
	hk, err := Secret(k).Hash()
	return HashedKey(hk), err
}

type HashedKey HashedSecret

func (h HashedKey) Compare(k Key) bool {
	return HashedSecret(h).Compare(Secret(k))
}

func TestSecretType(t *testing.T) {
	var k Key = "some key"
	var hk HashedKey
	hk, err := k.Hash()
	if err != nil {
		t.Fatalf("Key.Hash should not return an error, but got %v", err)
	}
	if hk == "" {
		t.Errorf("HashedKey should not be empty")
	}
	if !hk.Compare(k) {
		t.Errorf("HashedKey.Compare should return true")
	}
}
