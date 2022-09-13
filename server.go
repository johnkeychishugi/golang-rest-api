package main

import (
	"github.com/gin-gonic/gin"
	"github.com/johnkeychishugi/golang-api/config"
	"github.com/johnkeychishugi/golang-api/controllers"
	"github.com/johnkeychishugi/golang-api/middlewares"
	"github.com/johnkeychishugi/golang-api/repository"
	"github.com/johnkeychishugi/golang-api/services"
	"gorm.io/gorm"
)

var (
	//Database
	db *gorm.DB = config.SetUpDatabaseConnection()

	//Repositories
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)

	//Services
	jwtService  services.JWTService  = services.NewJWTService()
	userService services.UserService = services.NewUserService(userRepository)
	authService services.AuthService = services.NewAuthService(userRepository)
	bookService services.BookService = services.NewBookService(bookRepository)

	//Conrollers
	authController controllers.AuthController = controllers.NewAuthController(authService, jwtService)
	userController controllers.UserController = controllers.NewUserController(userService, jwtService)
	bookController controllers.BookController = controllers.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	server := gin.Default()

	authRoutes := server.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := server.Group("/api/user", middlewares.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	bookRoutes := server.Group("/api/books", middlewares.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.All)
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)

	}

	server.Run()
}
