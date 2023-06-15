package tests

import (
	"clothing-store-auth/internal/data"
	"clothing-store-auth/internal/validator"
	"testing"
	"time"
)

var user data.User

func TestValidateUser(t *testing.T) {
	user = data.User{
		ID:        12,
		Name:      "John",
		Money:     123,
		Email:     "john@gmail.com",
		Activated: true,
		Version:   1,
	}
	v := validator.New()
	user.Password.Set("1234")
	if data.ValidateUser(v, &user); v.Valid() {
		t.Fatalf(`Expected to be not valid but got: %v`, v.Valid())
	}

	v1 := validator.New()
	user.Password.Set("12345678")

	if data.ValidateUser(v1, &user); !v1.Valid() {
		t.Fatalf(`Expected to be valid but got: %v`, v1.Valid())
	}

}

func TestGenerateToken(t *testing.T) {
	token, err := data.GenerateToken(15, time.Hour, "test")
	if err != nil {
		t.Fatalf(`Expected to not get error but got: %v`, err)
	}
	if token.Plaintext == "" {
		t.Fatalf("Token plain text expected to be not empty, but it was empty")
	}
	if len(token.Hash) == 0 {
		t.Fatalf("Token hash expected to be not empty, but it was empty")
	}
	if token.Scope != "test" {
		t.Fatalf("Scope of the token is not the same with expected one")
	}
}

func TestValidateToken(t *testing.T) {
	v := validator.New()
	if data.ValidateTokenPlaintext(v, "1234"); v.Valid() {
		t.Fatalf(`Expected to be not valid but got: %v`, v.Valid())
	}
	v1 := validator.New()
	if data.ValidatePasswordPlaintext(v1, "DAHAULZJZPGLM7JMLDGZLGB4I4"); !v1.Valid() {
		t.Fatalf(`Expected to be valid but got: %v`, v1.Valid())
	}
}

func TestIncludeRole(t *testing.T) {
	role := data.Roles{"user", "admin"}
	value := role.Include("user")
	if value != true {
		t.Fatalf(`Expected value to be true, but got: %v`, value)
	}
	value = role.Include("moderator")
	if value != false {
		t.Fatalf(`Expected value to be false, but got: %v`, value)
	}

}
