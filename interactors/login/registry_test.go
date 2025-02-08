package login

import (
	"errors"
	"github.com/rogelioConsejo/go-regal/entities/user"
	"github.com/rogelioConsejo/go-regal/interactors/message"
	"strings"
	"testing"
)

func TestNewUserRegistry(t *testing.T) {
	t.Parallel()
	var r UserRegistry = NewUserRegistry(getSpyPersistence(), getSpyEmailClient())
	if r == nil {
		t.Fatal("user registry is nil")
	}
}

func TestUserRegistry_CreateUser(t *testing.T) {
	t.Parallel()
	persistence := getSpyPersistence()
	emailClient := getSpyEmailClient()
	registry := NewUserRegistry(persistence, emailClient)
	var userName user.Name = "testelio"
	var email user.Email = "testelio@emailprovider.com"
	u, err := user.New(userName, email)
	if err != nil {
		t.Fatal("unexpected error when creating user entity: ", err)
	}
	t.Run("it should create a user", func(t *testing.T) {
		_ = registry.CreateUser(u)
		exists, _ := registry.UserExists(u.Name())
		if !exists {
			t.Fatal("user was not created")
		}
	})
	t.Run("it should save the user through persistence", func(t *testing.T) {
		_ = registry.CreateUser(u)
		if len(persistence.users) == 0 {
			t.Fatal("no user was saved")
		}
	})
	t.Run("it should return an error if the user name already exists", func(t *testing.T) {
		err = registry.CreateUser(u)
		if err == nil {
			t.Fatal("expected error when creating repeated user")
		}
	})
	t.Run("it should return an error if persistence fails", func(t *testing.T) {
		persistence.setFailOnSaveUser(true)
		err = registry.CreateUser(u)
		if err == nil {
			t.Fatal("expected error when saving user")
		}
	})
	t.Run("it should send a confirmation email and save the confirmation code for the user", func(t *testing.T) {
		_ = registry.CreateUser(u)

		if len(emailClient.calls["Send"]) == 0 {
			t.Fatal("email was not sent")
		}
		sentEmail := emailClient.calls["Send"][1].(message.Message)

		code, ok := persistence.confirmationCodes[u.Name()]
		if !ok {
			t.Fatal("confirmation code was not saved")
		}
		if code == "" {
			t.Fatal("confirmation code is empty")
		}

		if !strings.Contains(sentEmail.Body, string(code)) {
			t.Fatal("confirmation code was not sent in the email")
		}
	})
}

func TestUserRegistry_UserExists(t *testing.T) {
	t.Parallel()
	persistence := getSpyPersistence()
	registry := NewUserRegistry(persistence, getSpyEmailClient())
	var userName user.Name = "testelio"
	var email user.Email = "testelio@email.com"
	u, err := user.New(userName, email)
	if err != nil {
		t.Fatal("unexpected error when creating user entity: ", err)
	}
	exists, _ := registry.UserExists(u.Name())
	if exists {
		t.Fatal("user should not exist yet")
	}
	_ = registry.CreateUser(u)
	exists, _ = registry.UserExists(u.Name())
	if !exists {
		t.Fatal("user should exist")
	}
	t.Run("it should check if the user was saved", func(t *testing.T) {
		if _, ok := persistence.users[u.Name()]; !ok {
			t.Fatal("user was not saved")
		}
		if persistence.calls["UserWasSaved"] == nil {
			t.Fatal("UserWasSaved method was not called")
		}
	})
	t.Run("it should return false if the user does not exist", func(t *testing.T) {
		exists, _ := registry.UserExists("nonexistent")
		if exists {
			t.Fatal("user should not exist")
		}
	})
	t.Run("it should return an error if the persistence fails", func(t *testing.T) {
		persistence.setFailOnCheckUser(true)
		exists, checkErr := registry.UserExists(u.Name())
		if checkErr == nil {
			t.Fatal("expected error when checking user")
		}
		if exists {
			t.Fatal("user should not exist")
		}
	})
}

func TestUserRegistry_ConfirmUserEmail(t *testing.T) {
	t.Run("it should confirm the user email using the confirmation code", func(t *testing.T) {
		// Setup
		persistence := getSpyPersistence()
		emailClient := getSpyEmailClient()
		registry := NewUserRegistry(persistence, emailClient)

		// Create a test user
		userName := user.Name("testelio")
		email := user.Email("testelio@emailprovider.com")
		u, err := user.New(userName, email)
		if err != nil {
			t.Fatal("unexpected error when creating user entity:", err)
		}

		// Create the user which should generate and save a confirmation code
		err = registry.CreateUser(u)
		if err != nil {
			t.Fatal("unexpected error when creating user:", err)
		}

		// Get the confirmation code that was saved
		code, ok := persistence.confirmationCodes[userName]
		if !ok {
			t.Fatal("confirmation code was not saved")
		}

		// Confirm the email
		err = registry.ConfirmUserEmail(userName, code)
		if err != nil {
			t.Fatal("unexpected error when confirming email:", err)
		}

		// Verify the email is confirmed
		confirmed, err := registry.UserEmailIsConfirmed(userName)
		if err != nil {
			t.Fatal("unexpected error when checking email confirmation:", err)
		}
		if !confirmed {
			t.Fatal("email should be confirmed but isn't")
		}
	})
	t.Run("it should return an error if the confirmation code is invalid", func(t *testing.T) {
		// Setup
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		// Create a test user
		userName := user.Name("testelio")
		email := user.Email("testelio@emailprovider.com")
		u, err := user.New(userName, email)
		if err != nil {
			t.Fatal("unexpected error when creating user entity:", err)
		}

		// Create the user first
		err = registry.CreateUser(u)
		if err != nil {
			t.Fatal("unexpected error when creating user:", err)
		}

		// Try to confirm with invalid code
		invalidCode := ConfirmationCode("wrong-code")
		err = registry.ConfirmUserEmail(userName, invalidCode)
		if err == nil {
			t.Fatal("expected error with invalid confirmation code but got none")
		}

		// Verify email remains unconfirmed
		confirmed, err := registry.UserEmailIsConfirmed(userName)
		if err != nil {
			t.Fatal("unexpected error when checking confirmation:", err)
		}
		if confirmed {
			t.Fatal("email should not be confirmed with invalid code")
		}
	})

	t.Run("it should return an error if the user does not exist", func(t *testing.T) {
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		nonExistentUser := user.Name("nonexistent")
		randomCode := ConfirmationCode("some-code")

		err := registry.ConfirmUserEmail(nonExistentUser, randomCode)
		if err == nil {
			t.Fatal("expected error with non-existent user but got none")
		}
	})
}

func TestUserRegistry_UserEmailIsConfirmed(t *testing.T) {
	t.Run("it should return true if the user email is confirmed", func(t *testing.T) {
		// Setup
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		// Create and setup test user
		userName := user.Name("testelio")
		email := user.Email("testelio@emailprovider.com")
		u, err := user.New(userName, email)
		if err != nil {
			t.Fatal("unexpected error when creating user entity:", err)
		}

		// Create user and confirm email
		err = registry.CreateUser(u)
		if err != nil {
			t.Fatal("unexpected error when creating user:", err)
		}

		code := persistence.confirmationCodes[userName]
		err = registry.ConfirmUserEmail(userName, code)
		if err != nil {
			t.Fatal("unexpected error when confirming email:", err)
		}

		// Check confirmation status
		confirmed, err := registry.UserEmailIsConfirmed(userName)
		if err != nil {
			t.Fatal("unexpected error when checking confirmation status:", err)
		}
		if !confirmed {
			t.Fatal("email should be confirmed")
		}
	})

	t.Run("it should return false if the user email is not confirmed", func(t *testing.T) {
		// Setup
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		// Create test user
		userName := user.Name("testelio")
		email := user.Email("testelio@emailprovider.com")
		u, err := user.New(userName, email)
		if err != nil {
			t.Fatal("unexpected error when creating user entity:", err)
		}

		// Create user but don't confirm email
		err = registry.CreateUser(u)
		if err != nil {
			t.Fatal("unexpected error when creating user:", err)
		}

		// Check confirmation status
		confirmed, err := registry.UserEmailIsConfirmed(userName)
		if err != nil {
			t.Fatal("unexpected error when checking confirmation status:", err)
		}
		if confirmed {
			t.Fatal("email should not be confirmed")
		}
	})
}

func TestUserRegistry_GetUserEmail(t *testing.T) {
	t.Run("it should return the user email", func(t *testing.T) {
		// Setup
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		// Create a test user
		userName := user.Name("testelio")
		expectedEmail := user.Email("testelio@emailprovider.com")
		u, err := user.New(userName, expectedEmail)
		if err != nil {
			t.Fatal("unexpected error when creating user entity:", err)
		}

		// Create and confirm user email
		err = registry.CreateUser(u)
		if err != nil {
			t.Fatal("unexpected error when creating user:", err)
		}

		// Confirm email
		code := persistence.confirmationCodes[userName]
		err = registry.ConfirmUserEmail(userName, code)
		if err != nil {
			t.Fatal("unexpected error when confirming email:", err)
		}

		// Get user email
		email, err := registry.GetUserEmail(userName)
		if err != nil {
			t.Fatal("unexpected error when getting user email:", err)
		}

		if email != expectedEmail {
			t.Fatalf("expected email %s but got %s", expectedEmail, email)
		}
	})
	t.Run("it return ErrGettingUserEmail when persistence fails", func(t *testing.T) {
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		persistence.setFailOnGetEmail(true)

		_, err := registry.GetUserEmail("test-user")
		if !errors.Is(err, ErrGettingUserEmail) {
			t.Fatalf("expected error to be of type ErrGettingUserEmail, got: %v", err)
		}
	})
	t.Run("it should return an error if the user does not exist", func(t *testing.T) {
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		_, err := registry.GetUserEmail("nonexistent")
		if err == nil {
			t.Fatal("expected error with non-existent user but got none")
		}

	})
	t.Run("it should return an error if the user email is not confirmed", func(t *testing.T) {
		persistence := getSpyPersistence()
		registry := NewUserRegistry(persistence, getSpyEmailClient())

		// Create a test user
		userName := user.Name("testelio")
		email := user.Email("testelio@emailprovider.com")
		u, err := user.New(userName, email)
		if err != nil {
			t.Fatal("unexpected error creating user:", err)
		}

		// Save user but don't confirm email
		err = registry.CreateUser(u)
		if err != nil {
			t.Fatal("unexpected error saving user:", err)
		}

		// Try to get email of unconfirmed user
		_, err = registry.GetUserEmail(userName)
		if !errors.Is(err, ErrEmailNotConfirmed) {
			t.Fatalf("expected %v but got: %v", ErrEmailNotConfirmed, err)
		}
		if !errors.Is(err, ErrGettingUserEmail) {
			t.Fatalf("expected %v but got: %v", ErrGettingUserEmail, err)
		}
	})
}

func getSpyPersistence() *spyPersistence {
	return &spyPersistence{
		users:             make(map[user.Name]user.User),
		calls:             make(map[string][]interface{}),
		confirmationCodes: make(map[user.Name]ConfirmationCode),
		confirmedEmails:   make(map[user.Name]bool),
	}
}

type spyPersistence struct {
	users             map[user.Name]user.User
	confirmationCodes map[user.Name]ConfirmationCode
	confirmedEmails   map[user.Name]bool // New field
	calls             map[string][]interface{}
	failOnSave        bool
	failOnCheck       bool
	failOnGetEmail    bool
}

func (s *spyPersistence) GetUserEmail(name user.Name) (user.Email, error) {
	if s.failOnGetEmail {
		return "", errors.New("failed to get email")
	}
	u, exists := s.users[name]
	if !exists {
		return "", errors.New("user not found")
	}
	return u.Email(), nil
}

func (s *spyPersistence) GetConfirmationCode(u user.Name) (ConfirmationCode, error) {
	code, exists := s.confirmationCodes[u]
	if !exists {
		return "", errors.New("confirmation code not found")
	}
	return code, nil
}

func (s *spyPersistence) MarkEmailAsConfirmed(u user.Name) error {
	s.confirmedEmails[u] = true
	return nil
}

func (s *spyPersistence) IsEmailConfirmed(u user.Name) (bool, error) {
	confirmed, exists := s.confirmedEmails[u]
	if !exists {
		return false, nil
	}
	return confirmed, nil
}

func (s *spyPersistence) UserWasSaved(u user.Name) (bool, error) {
	if s.failOnCheck {
		return false, errors.New("failed to check user")
	}
	s.calls["UserWasSaved"] = append(s.calls["UserWasSaved"], u)
	_, ok := s.users[u]
	return ok, nil
}

func (s *spyPersistence) SaveUser(u user.User) error {
	if s.failOnSave {
		return errors.New("failed to save user")
	}
	s.calls["SaveUser"] = append(s.calls["SaveUser"], u)
	s.users[u.Name()] = u
	return nil
}

func (s *spyPersistence) SaveConfirmationCode(u user.Name, c ConfirmationCode) error {
	s.calls["SaveConfirmationCode"] = append(s.calls["SaveConfirmationCode"], u, c)
	s.confirmationCodes[u] = c
	return nil
}

func (s *spyPersistence) setFailOnSaveUser(sw bool) {
	s.failOnSave = sw
}

func (s *spyPersistence) setFailOnCheckUser(b bool) {
	s.failOnCheck = b
}

func (s *spyPersistence) setFailOnGetEmail(b bool) {
	s.failOnGetEmail = b
}

var _ RegistryPersistence = &spyPersistence{}

type spyEmailClient struct {
	calls map[string][]interface{}
}

func (s *spyEmailClient) Send(address message.Address, m message.Message) error {
	s.calls["Send"] = append(s.calls["Send"], address, m)
	return nil
}

func getSpyEmailClient() *spyEmailClient {
	return &spyEmailClient{
		calls: make(map[string][]interface{}),
	}
}
