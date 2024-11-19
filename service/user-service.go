package service

import "auth/presentation/dto"

type UserService interface {
	Register(registerRequest dto.RegisterRequest) error
}
