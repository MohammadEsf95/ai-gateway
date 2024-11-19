package repository

import "auth/presentation/dto"

type UserRepository interface {
	Register(registerRequest dto.RegisterRequest) error
}
