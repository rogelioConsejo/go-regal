package user

import "testing"

func TestNew(t *testing.T) {
	t.Run("it should return a user with the correct name and email", func(t *testing.T) {
		t.Parallel()
		var userName Name = "MrAnderson"
		var email Email = "mranderson@email.com"
		u, err := New(userName, email)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if u == nil {
			t.Error("User is nil")
		}
		t.Run("it should return a user with the correct name", func(t *testing.T) {
			if u.Name() != userName {
				t.Errorf("Expected: %v, Got: %v", userName, u.Name())
			}
		})
		t.Run("it should return a user with the correct email", func(t *testing.T) {
			if u.Email() != email {
				t.Errorf("Expected: %v, Got: %v", email, u.Email())
			}
		})
	})
	t.Run("it should return an error if the name is empty", func(t *testing.T) {
		t.Parallel()
		var userName Name = ""
		var email Email = "someemail@email.com"
		u, err := New(userName, email)
		if err == nil {
			t.Error("Error is nil")
		}
		if u != nil {
			t.Error("User is not nil")
		}
	})
	t.Run("it should return an error if the email is empty or invalid", func(t *testing.T) {
		t.Parallel()
		t.Run("it should return an error if the email is empty", func(t *testing.T) {
			t.Parallel()
			var userName Name = "MrAnderson"
			var email Email = ""
			u, err := New(userName, email)
			if err == nil {
				t.Error("Error is nil")
			}
			if u != nil {
				t.Error("User is not nil")
			}
		})
		t.Run("it should return an error if the email has no @ symbol", func(t *testing.T) {
			t.Parallel()
			var userName Name = "MrAnderson"
			var email Email = "mrandersonemail.com"
			u, err := New(userName, email)
			if err == nil {
				t.Error("Error is nil")
			}
			if u != nil {
				t.Error("User is not nil")
			}
		})
		t.Run("it should return an error if the @ symbol is at the beginning of the email", func(t *testing.T) {
			t.Parallel()
			var userName Name = "MrAnderson"
			var email Email = "@mrandersonemail.com"
			u, err := New(userName, email)
			if err == nil {
				t.Error("Error is nil")
			}
			if u != nil {
				t.Error("User is not nil")
			}
		})
		t.Run("it should return an error if the email has no . symbol", func(t *testing.T) {
			t.Parallel()
			var userName Name = "MrAnderson"
			var email Email = "mranderson@"
			u, err := New(userName, email)
			if err == nil {
				t.Error("Error is nil")
			}
			if u != nil {
				t.Error("User is not nil")
			}
		})
		t.Run("it should return an error if the @ symbol is after the last . symbol", func(t *testing.T) {
			t.Parallel()
			var userName Name = "MrAnderson"
			var email Email = "mrandersonemail.com@"
			u, err := New(userName, email)
			if err == nil {
				t.Error("Error is nil")
			}
			if u != nil {
				t.Error("User is not nil")
			}
		})
		t.Run("it should return an error if the email ends with a . symbol", func(t *testing.T) {
			t.Parallel()
			var userName Name = "MrAnderson"
			var email Email = "mranderson@email.com."
			u, err := New(userName, email)
			if err == nil {
				t.Error("Error is nil")
			}
			if u != nil {
				t.Error("User is not nil")
			}
		})
	})
}

var _ User = user{}
