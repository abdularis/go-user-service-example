package handler

import (
	"github.com/gin-gonic/gin"
	"go-user-service-example/app/core"
	"go-user-service-example/app/middleware"
	"go-user-service-example/app/utils"
	"gorm.io/gorm"
	"net/http"
)

type UserHandler struct {
	userRepo core.UserRepository
}

func NewUserHandler(repository core.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repository}
}

func (h *UserHandler) CreateNewUser(c *gin.Context) {
	createRequest := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
		return
	}

	user, err := h.userRepo.FindUserByUserName(c.Request.Context(), createRequest.UserName)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	if user != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Message: "User already exists"})
		return
	}

	hashedPassword, err := core.HashPassword(createRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	user = &core.User{
		UserName: createRequest.UserName,
		Password: hashedPassword,
		Role:     core.DefaultUser,
	}
	if err := h.userRepo.SaveUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, utils.Response{Message: "success", Data: user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := utils.StrToUint(c.Param("userID"))
	updateRequest := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Message: err.Error()})
		return
	}

	user, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	existingUser, err := h.userRepo.FindUserByUserName(c.Request.Context(), updateRequest.UserName)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	if existingUser != nil && existingUser.ID != user.ID {
		c.JSON(http.StatusBadRequest, utils.Response{Message: "Username already used by other user"})
		return
	}

	hashedPassword, err := core.HashPassword(updateRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	user.UserName = updateRequest.UserName
	user.Password = hashedPassword
	if err := h.userRepo.SaveUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Message: "update user success", Data: user})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := utils.StrToUint(c.Param("userID"))
	user, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Message: "ok", Data: user})
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	users, err := h.userRepo.FindAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Message: "ok", Data: users})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusBadRequest, utils.Response{Message: "No authenticated user found"})
		return
	}
	// just return the user from auth cause its already complete user object from token
	c.JSON(http.StatusOK, utils.Response{Message: "success", Data: user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := utils.StrToUint(c.Param("userID"))
	if err := h.userRepo.DeleteUserByID(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Message: "success"})
}
