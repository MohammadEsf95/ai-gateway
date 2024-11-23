package repository

import (
	"auth/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	FindById(id uint) (*domain.User, error)
	DeleteById(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindById(id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) DeleteById(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
