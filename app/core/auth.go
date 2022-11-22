package core

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"strconv"
	"time"
)

type TokenType string

var (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

var tokenSecrets = map[TokenType][]byte{}

func SetJwtSecret(secretAccessToken string, secretRefreshToken string) {
	tokenSecrets[AccessToken] = []byte(secretAccessToken)
	tokenSecrets[RefreshToken] = []byte(secretRefreshToken)
}

func GenerateJWTByUser(typ TokenType, user User) (string, error) {
	tokenBuilder := jwt.NewBuilder().
		Claim("uid", fmt.Sprintf("%d", user.ID)).
		Claim("role", user.Role).
		Subject(user.UserName).
		IssuedAt(time.Now())
	if typ == AccessToken {
		tokenBuilder.Expiration(time.Now().Add(time.Hour * 24)) // expire on 1 day
	}

	t, err := tokenBuilder.Build()
	if err != nil {
		return "", err
	}

	// use HS256 algorithm for signing to simplify code by allowing us to use key using byte string
	signed, err := jwt.Sign(t, jwa.HS256, tokenSecrets[typ])
	if err != nil {
		return "", err
	}
	return string(signed), nil
}

func VerifyJWT(typ TokenType, token string) (*User, error) {
	t, err := jwt.Parse([]byte(token), jwt.WithVerify(jwa.HS256, tokenSecrets[typ]))
	if err != nil {
		return nil, err
	}

	err = jwt.Validate(t)
	if err != nil {
		return nil, err
	}

	u := User{
		UserName: t.Subject(),
	}

	if uid, ok := t.Get("uid"); ok {
		userId, err := strconv.Atoi(uid.(string))
		if err != nil {
			return nil, err
		}
		u.ID = uint(userId)
	}

	if role, ok := t.Get("role"); ok {
		u.SetRoleFromStr(role.(string))
	}
	return &u, nil
}
