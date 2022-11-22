package handler

import (
	"github.com/gin-gonic/gin"
	"go-user-service-example/app/core"
	"go-user-service-example/app/middleware"
	"go-user-service-example/app/utils"
	"net/http"
)

type AuthHandler struct {
	userRepo core.UserRepository
}

func NewAuthHandler(repository core.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: repository}
}

func (h *AuthHandler) Login(c *gin.Context) {
	loginRequest := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
		return
	}

	user, err := h.userRepo.FindUserByUserName(c.Request.Context(), loginRequest.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	err = user.VerifyPassword(loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Message: "Failed to verify password"})
		return
	}

	token, err := core.GenerateJWTByUser(core.AccessToken, *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	refreshToken, err := core.GenerateJWTByUser(core.RefreshToken, *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Message: "success",
		Data: map[string]interface{}{
			"accessToken":  token,
			"refreshToken": refreshToken,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	token, err := core.GenerateJWTByUser(core.AccessToken, *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Message: "success",
		Data: map[string]interface{}{
			"accessToken": token,
		},
	})
}

//func (h *AuthHandler) Register(c *gin.Context) {
//	registerRequest := struct {
//		UserName string `json:"userName"`
//		Password string `json:"password"`
//	}{}
//	if err := c.Bind(&registerRequest); err != nil {
//		c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
//		return
//	}
//
//	user, err := h.userRepo.FindUserByUserName(c.Request.Context(), registerRequest.UserName)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
//		return
//	}
//
//	if user != nil {
//		c.JSON(http.StatusBadRequest, utils.Response{Message: "User already exists"})
//		return
//	}
//
//	hashedPassword, err := core.HashPassword(registerRequest.Password)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
//		return
//	}
//
//	user = &core.User{
//		UserName: registerRequest.UserName,
//		Password: hashedPassword,
//		Role:     core.DefaultUser,
//	}
//	if err := h.userRepo.SaveUser(c.Request.Context(), user); err != nil {
//		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
//		return
//	}
//	c.JSON(http.StatusCreated, utils.Response{Message: "success"})
//}
