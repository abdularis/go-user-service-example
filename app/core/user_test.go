package core

import "testing"

func TestHashVerifyPassword(t *testing.T) {
	hashed, err := HashPassword("admin")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hashed)

	user := User{Password: hashed}
	err = user.VerifyPassword("admin")
	if err != nil {
		t.Fatal(err)
	}
}
