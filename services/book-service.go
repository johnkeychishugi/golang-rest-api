package services

import (
	"fmt"
	"log"

	"github.com/johnkeychishugi/golang-api/models"
	"github.com/johnkeychishugi/golang-api/repository"
	"github.com/johnkeychishugi/golang-api/validations"
	"github.com/mashingan/smapping"
)

type BookService interface {
	Insert(b validations.BookCreateValidation) models.Book
	Update(b validations.BookUpdateValidation) models.Book
	Delete(b models.Book)
	All() []models.Book
	FindByID(bookID uint64) models.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (service *bookService) Insert(b validations.BookCreateValidation) models.Book {
	book := models.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}

	res := service.bookRepository.InsertBook(book)
	return res
}

func (service *bookService) Update(b validations.BookUpdateValidation) models.Book {
	book := models.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.bookRepository.UpdateBook(book)
	return res
}

func (service *bookService) Delete(b models.Book) {
	service.bookRepository.DeleteBook(b)
}

func (service *bookService) All() []models.Book {
	return service.bookRepository.AllBook()
}

func (service *bookService) FindByID(bookID uint64) models.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	b := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", b.UserID)

	return userID == id
}
