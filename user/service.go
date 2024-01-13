package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	RegisterInput(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsAvailableEmail(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, fileLocation string) (User, error)
	GetUserByID(id int) (User, error)
}

type usecase struct {
	repository Repository
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{repository}
}

func (u *usecase) RegisterInput(input RegisterUserInput) (User, error) {
	user := User{}

	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email

	pass, errPass := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if errPass != nil {
		return user, errPass
	}
	user.PasswordHash = string(pass)

	result, err := u.repository.Save(user)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *usecase) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := u.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("user not found")
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if errPass != nil {
		return user, errPass
	}

	return user, nil
}

func (u *usecase) IsAvailableEmail(input CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := u.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (u *usecase) SaveAvatar(id int, fileLocation string) (User, error) {
	user, err := u.repository.FindById(id)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	uploadAvtr, errUpAvtr := u.repository.Update(user)
	if errUpAvtr != nil {
		return uploadAvtr, errUpAvtr
	}

	return uploadAvtr, nil
}

func (u *usecase) GetUserByID(id int) (User, error) {
	user, err := u.repository.FindById(id)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on with that id")
	}

	return user, nil
}
