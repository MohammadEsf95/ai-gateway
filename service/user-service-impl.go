package service

import (
	"auth/presentation/dto"
	"auth/repository"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{userRepository: userRepository}
}

func (s *UserServiceImpl) Register(registerRequest dto.RegisterRequest) error {
	return s.userRepository.Register(registerRequest)
}
