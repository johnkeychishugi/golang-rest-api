package services

import (
	"log"

	"github.com/johnkeychishugi/golang-api/models"
	"github.com/johnkeychishugi/golang-api/repository"
	"github.com/johnkeychishugi/golang-api/validations"
	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user validations.UserUpdateValidation) models.User
	Profile(userID string) models.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user validations.UserUpdateValidation) models.User {
	userToUpdate := models.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) models.User {
	return service.userRepository.ProfileUser(userID)
}
