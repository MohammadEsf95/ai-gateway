package repository

import (
	"auth/entities"
	"auth/infrastructure"
	"auth/presentation/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Register(registerRequest dto.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return infrastructure.ErrFaildHashPassword
	}
	return r.db.Create(&entities.User{FirstName: registerRequest.FirstName, LastName: registerRequest.LastName, Email: registerRequest.Email, Password: string(hashedPassword)}).Error
}
