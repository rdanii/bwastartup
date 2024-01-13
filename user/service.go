package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	RegisterInput(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
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
