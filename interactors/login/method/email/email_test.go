package email

import (
	"github.com/rogelioConsejo/go-regal/entities/user"
	"github.com/rogelioConsejo/go-regal/interactors/login"
	"testing"
	"time"
)

func TestGetEmailMethod(t *testing.T) {
	t.Parallel()
	var m login.AuthenticationMethod = GetEmailMethod(getSpyPersistence(), getSpyEmailClient())
	if m == nil {
		t.Error("GetEmailMethod should return an authentication method")
	}
	t.Run("it should return an emailMethod instance", func(t *testing.T) {
		_, ok := m.(*emailMethod)
		if !ok {
			t.Error("expected an emailMethod instance")
		}
		t.Run("it should have a default 15 minutes expiration time", func(t *testing.T) {
			if m.(*emailMethod).expiration != 15*time.Minute {
				t.Errorf("expected expiration time to be 15, got %v", m.(*emailMethod).expiration)
			}
		})
	})
}

func TestEmailMethod_SetupAuthenticationAttempt(t *testing.T) {
	t.Parallel()
	var m login.AuthenticationMethod = GetEmailMethod(getSpyPersistence(), getSpyEmailClient())
	m.(*emailMethod).persistence = getSpyPersistence()
	m.(*emailMethod).emailSender = getSpyEmailClient()
	u, err := user.New("test", "test@mail.com")
	if err != nil {
		t.Fatalf("error creating user entity: %v", err)
	}
	t.Run("it should persist a new token for each attempt", func(t *testing.T) {
		oldToken := ""
		t.Run("it should generate a login token, associate it with the user and persist it", func(t *testing.T) {
			err = m.SetupAuthenticationAttempt(u)
			if err != nil {
				t.Fatalf("unexpected error %s", err.Error())
			}
			savedTokens := m.(*emailMethod).persistence.(*spyPersistence).savedLoginTokens
			if len(savedTokens) != 1 {
				t.Errorf("expected to save 1 token, but got %v", len(savedTokens))
			}
			oldToken = string(savedTokens[u.Name()].HashedToken)
		})
		t.Run("it should generate a different login token for each authentication attempt", func(t *testing.T) {
			err = m.SetupAuthenticationAttempt(u)
			if err != nil {
				t.Fatalf("unexpected error %s", err.Error())
			}
			savedTokens := m.(*emailMethod).persistence.(*spyPersistence).savedLoginTokens
			if len(savedTokens) != 1 {
				t.Errorf("expected to save 1 token, but got %v", len(savedTokens))
			}
			newToken := string(savedTokens[u.Name()].HashedToken)
			if oldToken == newToken {
				t.Error("expected to save a different token")
			}
		})
	})
	t.Run("it should send the login token to the user's email using an email client", func(t *testing.T) {
		sentEmails := m.(*emailMethod).emailSender.(*spyEmailClient).sent
		if len(sentEmails[u.Email()]) < 1 {
			t.Error("no emails have been sent")
		}

	})

}

func TestEmailMethod_Authenticate(t *testing.T) {
	t.Parallel()
	var m login.AuthenticationMethod = GetEmailMethod(getSpyPersistence(), getSpyEmailClient())
	u, err := user.New("test", "test@mail.com")
	if err != nil {
		t.Fatalf("error creating user entity: %v", err)
	}
	err = m.SetupAuthenticationAttempt(u)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	t.Run("it should return true if the correct token is used", func(t *testing.T) {
		tokens := m.(*emailMethod).emailSender.(*spyEmailClient).sent[u.Email()]
		lastToken := tokens[len(tokens)-1]
		authenticated, authErr := m.Authenticate(u, login.Credential(lastToken))
		if authErr != nil {
			t.Fatalf("error authenticating user: %s", authErr.Error())
		}
		if !authenticated {
			t.Error("expected user to be authenticated")
		}
	})
	t.Run("it should return false if a wrong token is used", func(t *testing.T) {
		authenticated, authErr := m.Authenticate(u, "some-token")
		if authErr != nil {
			t.Fatalf("error authenticating user: %s", authErr.Error())
		}
		if authenticated {
			t.Error("expected user not to be authenticated")
		}
	})
	t.Run("it should return false if the user did not request a token", func(t *testing.T) {
		otherUser, _ := user.New("other", "other@email.com")
		authenticated, _ := m.Authenticate(otherUser, "some-token")
		if authenticated {
			t.Error("expected user not to be authenticated")
		}
	})
	t.Run("it should return false if the token is expired", func(t *testing.T) {
	})
}

func getSpyEmailClient() Sender {
	return &spyEmailClient{
		sent: make(map[user.Email][]Token),
	}
}

type spyEmailClient struct {
	sent map[user.Email][]Token
}

func (c *spyEmailClient) sendToken(email user.Email, tk Token) error {
	c.sent[email] = append(c.sent[email], tk)
	return nil
}

func getSpyPersistence() Persistence {
	return &spyPersistence{
		savedLoginTokens: make(map[user.Name]SavedToken),
	}
}

type spyPersistence struct {
	savedLoginTokens map[user.Name]SavedToken
}

type SavedToken struct {
	HashedToken
	Expiration time.Time
}

func (p *spyPersistence) GetLoginToken(n user.Name) (HashedToken, time.Time, error) {
	return p.savedLoginTokens[n].HashedToken, p.savedLoginTokens[n].Expiration, nil
}

func (p *spyPersistence) SaveLoginToken(u user.User, token HashedToken, exp time.Time) error {
	p.savedLoginTokens[u.Name()] = SavedToken{
		HashedToken: token,
		Expiration:  exp,
	}
	return nil
}
