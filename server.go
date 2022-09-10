package main

import (
	"github.com/gin-gonic/gin"
	"github.com/johnkeychishugi/golang-api/config"
	"github.com/johnkeychishugi/golang-api/controllers"
	"github.com/johnkeychishugi/golang-api/repository"
	"github.com/johnkeychishugi/golang-api/services"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                   = config.SetUpDatabaseConnection()
	userRepository repository.UserRepository  = repository.NewUserRepository(db)
	jwtService     services.JWTService        = services.NewJWTService()
	authService    services.AuthService       = services.NewAuthService(userRepository)
	authController controllers.AuthController = controllers.NewAuthController(authService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	server := gin.Default()

	authRoutes := server.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	server.Run()
}
