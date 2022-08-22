package main

import (
	"github.com/gin-gonic/gin"
	"github.com/johnkeychishugi/golang-api/config"
	controller "github.com/johnkeychishugi/golang-api/controllers"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetUpDatabaseConnection()
	authController controller.AuthController = controller.NewAuthController()
)

func main() {
	server := gin.Default()

	authRoutes := server.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	server.Run()
}
