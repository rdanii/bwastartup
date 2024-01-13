package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindById(id int) (User, error)
	Update(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	var input RegisterUserInput
	input.Email = user.Email

	err := r.db.Debug().First(&user, "email = ?", input.Email).Error
	if err == nil {
		return user, errors.New("email already exists")
	}

	err = r.db.Debug().Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User

	err := r.db.Debug().Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindById(id int) (User, error) {
	var user User
	err := r.db.Debug().Find(&user, "id = ?", id).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) Update(user User) (User, error) {
	err := r.db.Debug().Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
