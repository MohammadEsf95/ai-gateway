package service

import (
	"auth/domain"
	"auth/presentation/dto"
	"auth/repository"
	"errors"

	"auth/pkg/utils"

	"github.com/markbates/goth"
)

type UserService interface {
	Register(req dto.RegisterRequest) (*domain.User, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	UpsertGoogleUser(gothUser goth.User) (*dto.AuthResponse, error)
	GetUserById(id uint) (*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwtUtil  utils.JWTUtil
}

func NewUserService(userRepo repository.UserRepository, jwtUtil utils.JWTUtil) UserService {
	return &userService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *userService) Register(req dto.RegisterRequest) (*domain.User, error) {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Provider: "email",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtUtil.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *userService) UpsertGoogleUser(gothUser goth.User) (*dto.AuthResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(gothUser.Email)
	if err != nil {
		return nil, err
	}

	var user *domain.User

	if existingUser == nil {
		user = &domain.User{
			Email:      gothUser.Email,
			Name:       gothUser.Name,
			Provider:   "google",
			ProviderId: gothUser.UserID,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	} else {
		existingUser.Name = gothUser.Name
		existingUser.Provider = "google"
		existingUser.ProviderId = gothUser.UserID
		if err := s.userRepo.Update(existingUser); err != nil {
			return nil, err
		}
		user = existingUser
	}

	token, err := s.jwtUtil.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *userService) GetUserById(id uint) (*domain.User, error) {
	return s.userRepo.FindById(id)
}
