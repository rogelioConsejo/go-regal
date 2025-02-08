package login

import (
	"github.com/rogelioConsejo/go-regal/entities/user"
	"testing"
)

func TestNewAccess(t *testing.T) {
	t.Parallel()
	am := getAuthMethod()
	var a Access = NewAccess(am)
	if a == nil {
		t.Error("expected an access instance")
	}
}

func TestAccess_PrepareAuthenticationAttempt(t *testing.T) {
	t.Run("it should call the SetupAuthenticationAttempt method of the default authentication method", func(t *testing.T) {
		t.Parallel()
		am := getAuthMethod()
		a := NewAccess(am)
		u, _ := user.New("test", "test@email.com")
		err := a.PrepareAuthentication(u)
		if err != nil {
			t.Fatalf("unexpected error %s", err.Error())
		}
		if _, ok := am.calls["SetupAuthenticationAttempt"]; !ok {
			t.Error("expected SetupAuthenticationAttempt to be called")
		}
		if am.calls["SetupAuthenticationAttempt"][0].(user.User).Name() != u.Name() {
			t.Errorf("expected user to be %s, got %s", u.Name(), am.calls["SetupAuthenticationAttempt"][0].(user.User).Name())
		}
		if am.calls["SetupAuthenticationAttempt"][0].(user.User).Email() != u.Email() {
			t.Errorf("expected user to be %s, got %s", u.Email(), am.calls["SetupAuthenticationAttempt"][0].(user.User).Email())
		}
	})
}

func TestAccess_PerformAuthentication(t *testing.T) {
	t.Run("it should call the Authenticate method of the default authentication method", func(t *testing.T) {
		t.Parallel()
		am := getAuthMethod()
		a := NewAccess(am)
		u, _ := user.New("test", "test@mail.com")
		var cred Credential = "test"
		_, err := a.PerformAuthentication(u, cred)
		if err != nil {
			t.Fatalf("unexpected error %s", err.Error())
		}
		if _, ok := am.calls["Authenticate"]; !ok {
			t.Fatal("expected Authenticate to be called")
		}
		if am.calls["Authenticate"][0].(user.User).Name() != u.Name() {
			t.Errorf("expected user to be %s, got %s", u.Name(), am.calls["Authenticate"][0].(user.User).Name())
		}
		if am.calls["Authenticate"][0].(user.User).Email() != u.Email() {
			t.Errorf("expected user to be %s, got %s", u.Email(), am.calls["Authenticate"][0].(user.User).Email())
		}
		if am.calls["Authenticate"][1].(Credential) != cred {
			t.Errorf("expected credential to be %s, got %s", cred, am.calls["Authenticate"][1].(Credential))
		}
	})
}

func getAuthMethod() *spyAuthenticationMethod {
	return &spyAuthenticationMethod{
		calls: make(map[string][]interface{}),
	}
}

type spyAuthenticationMethod struct {
	calls map[string][]interface{}
}

func (m *spyAuthenticationMethod) SetupAuthenticationAttempt(u user.User) error {
	m.calls["SetupAuthenticationAttempt"] = append(m.calls["SetupAuthenticationAttempt"], u)
	return nil
}

func (m *spyAuthenticationMethod) Authenticate(u user.User, credential Credential) (bool, error) {
	m.calls["Authenticate"] = append(m.calls["Authenticate"], u, credential)
	return true, nil
}
