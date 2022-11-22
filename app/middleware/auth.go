package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"go-user-service-example/app/core"
	"go-user-service-example/app/utils"
	"log"
	"net/http"
	"strings"
)

func GetAuthUser(c *gin.Context) *core.User {
	obj, ok := c.Get("user")
	if ok {
		return obj.(*core.User)
	}
	return nil
}

func verifyJwtFromAuthorizationHeader(typ core.TokenType, c *gin.Context) *core.User {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, utils.Response{Message: "Please provide authorization token"})
		return nil
	}

	if strings.HasPrefix(strings.ToLower(token), "bearer") {
		token = strings.TrimSpace(token[len("bearer"):])
	}

	user, err := core.VerifyJWT(typ, token)
	if err != nil {
		if jwt.IsValidationError(err) {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Message: "Token validation error (possibly token expired)",
				Data:    "token-validation-error",
			})
			return nil
		}
		log.Printf("Failed to verify token: %s\n", err)
		c.JSON(http.StatusUnauthorized, utils.Response{Message: "Failed to verify token"})
		return nil
	}

	return user
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := verifyJwtFromAuthorizationHeader(core.AccessToken, c)
		if user == nil {
			c.Abort()
			return
		}

		if user.Role != core.Admin {
			c.JSON(http.StatusForbidden, utils.Response{Message: "You don't have access to this resources"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func AllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := verifyJwtFromAuthorizationHeader(core.AccessToken, c)
		if user == nil {
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := verifyJwtFromAuthorizationHeader(core.RefreshToken, c)
		if user == nil {
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
