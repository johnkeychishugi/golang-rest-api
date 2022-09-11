package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/johnkeychishugi/golang-api/helpers"
	"github.com/johnkeychishugi/golang-api/services"
	"github.com/johnkeychishugi/golang-api/validations"
)

type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
}

type userController struct {
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) Update(context *gin.Context) {
	var userUpdateValidation validations.UserUpdateValidation
	errValidation := context.ShouldBind(&userUpdateValidation)
	if errValidation != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errValidation.Error(), helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadGateway, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateValidation.ID = id
	u := c.userService.Update(userUpdateValidation)
	res := helpers.BuildResponse(true, "OK!", u)
	context.JSON(http.StatusOK, res)
}

func (c *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := c.userService.Profile(id)
	res := helpers.BuildResponse(true, "OK", user)
	context.JSON(http.StatusOK, res)
}
