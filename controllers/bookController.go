package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/johnkeychishugi/golang-api/helpers"
	"github.com/johnkeychishugi/golang-api/models"
	"github.com/johnkeychishugi/golang-api/services"
	"github.com/johnkeychishugi/golang-api/validations"
)

type BookController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type bookController struct {
	bookService services.BookService
	jwtService  services.JWTService
}

func NewBookController(bookServ services.BookService, jwtServ services.JWTService) BookController {
	return &bookController{
		bookService: bookServ,
		jwtService:  jwtServ,
	}
}

func (c *bookController) All(context *gin.Context) {
	var books []models.Book = c.bookService.All()
	res := helpers.BuildResponse(true, "OK!", books)
	context.JSON(http.StatusOK, res)
}

func (c *bookController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helpers.BuildErrorResponse("No param id was found", err.Error(), helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book models.Book = c.bookService.FindByID(id)
	if (book == models.Book{}) {
		res := helpers.BuildErrorResponse("Data not found", "No data with given id", helpers.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helpers.BuildResponse(true, "OK!", book)
		context.JSON(http.StatusOK, res)
	}
}

func (c *bookController) Insert(context *gin.Context) {
	var bookCreateValidation validations.BookCreateValidation
	errValidation := context.ShouldBind(&bookCreateValidation)
	if errValidation != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errValidation.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)

		if err == nil {
			bookCreateValidation.UserID = uint16(convertedUserID)
		}

		result := c.bookService.Insert(bookCreateValidation)
		response := helpers.BuildResponse(true, "OK!", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *bookController) Update(context *gin.Context) {
	var bookUpdateValidation validations.BookUpdateValidation
	errValidation := context.ShouldBind(&bookUpdateValidation)

	if errValidation != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errValidation.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	} else {
		authHeader := context.GetHeader("Authorization")
		token, errToken := c.jwtService.ValidateToken(authHeader)

		if errToken != nil {
			panic(errToken.Error())
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := fmt.Sprintf("%v", claims["user_id"])

		if c.bookService.IsAllowedToEdit(userID, bookUpdateValidation.ID) {
			id, errID := strconv.ParseUint(userID, 10, 64)

			if errID == nil {
				bookUpdateValidation.UserID = uint16(id)
			}

			result := c.bookService.Update(bookUpdateValidation)
			response := helpers.BuildResponse(true, "OK!", result)
			context.JSON(http.StatusOK, response)
		} else {
			response := helpers.BuildErrorResponse("You dont have the right permission", "You are not the owber", helpers.EmptyObj{})
			context.JSON(http.StatusForbidden, response)
		}
	}
}

func (c *bookController) Delete(context *gin.Context) {
	var book models.Book
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)

	if err != nil {
		response := helpers.BuildErrorResponse("Failed to get ID", "No param id were found", helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}

	book.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.Delete(book)
		response := helpers.BuildResponse(true, "Deleted", helpers.EmptyObj{})
		context.JSON(http.StatusOK, response)
	} else {
		response := helpers.BuildErrorResponse("You dont have the right permission", "You are not the owber", helpers.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}
