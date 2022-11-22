package core

import (
	"testing"
)

func TestGenerateVerify(t *testing.T) {
	SetJwtSecret("mysecretforjwtsigning", "mysecretforjwtsigning2")
	user1 := User{
		ID:       1,
		UserName: "aris",
		Password: "",
		Role:     "admin",
	}
	token, err := GenerateJWTByUser(AccessToken, user1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)

	u, err := VerifyJWT(AccessToken, token)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user1)
	t.Log(*u)

	if user1 != *u {
		t.Fatal("user result doesn't match")
	}
}

func TestVerify(t *testing.T) {
	SetJwtSecret("my-jwt-secret", "my-jwt-secret-2")
	u, err := VerifyJWT(AccessToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NjkwOTQyMDgsInJvbGUiOiJhZG1pbiIsInN1YiI6ImFkbWluIiwidWlkIjoiMSJ9.LPgfy6IBw9QM2gDtPp-9guAATkJ-vqX9jEVZY9C8tds")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(u)
}
