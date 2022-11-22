package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-user-service-example/app/config"
	"go-user-service-example/app/core"
	"go-user-service-example/app/handler"
	"go-user-service-example/app/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func openMysqlConn(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func main() {
	cfg := config.Get()
	core.SetJwtSecret(cfg.JWTSecretAccessToken, cfg.JWTSecretRefreshToken)

	db, err := openMysqlConn(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := core.NewUserRepositoryMysql(db)
	authHandler := handler.NewAuthHandler(userRepo)
	userHandler := handler.NewUserHandler(userRepo)

	r := gin.Default()
	// unauthenticated access by anyone
	r.POST("/auth/login", authHandler.Login)
	refreshApi := r.Group("/").Use(middleware.RefreshToken())
	refreshApi.GET("/auth/refreshToken", authHandler.RefreshToken)

	// can be accessed by all authenticated users
	allUsersApi := r.Group("/").Use(middleware.AllUsers())
	allUsersApi.GET("/userInfo", userHandler.GetUserInfo)

	// can only be accessed by authenticated admin
	adminOnlyApi := r.Group("/").Use(middleware.AdminOnly())
	adminOnlyApi.POST("/users", userHandler.CreateNewUser)
	adminOnlyApi.PUT("/users/:userID", userHandler.UpdateUser)
	adminOnlyApi.DELETE("/users/:userID", userHandler.DeleteUser)
	adminOnlyApi.GET("/users/:userID", userHandler.GetUserByID)
	adminOnlyApi.GET("/users", userHandler.GetUserList)

	if err := r.Run(fmt.Sprintf(":%d", cfg.HostPort)); err != nil {
		log.Fatal(err)
	}
}
