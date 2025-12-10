package service

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/qs-lzh/movie-reservation/internal/security"
)

type AuthService interface {
	Login(username, password, key string) (token string, err error)
	ValidateToken(token string) (claims jwt.MapClaims, err error)
}

// jwtAuthService relies on UserService
type jwtAuthService struct {
	userService UserService
}

var _ AuthService = (*jwtAuthService)(nil)

func NewJWTAuthService(userService UserService) *jwtAuthService {
	return &jwtAuthService{
		userService: userService,
	}
}

func (s *jwtAuthService) Login(username, password, key string) (token string, err error) {

	isValid, err := s.userService.ValidateUser(username, password)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", ErrInvalidCredential
	}
	userRole, err := s.userService.GetUserRoleByName(username)
	if err != nil {
		return "", err
	}
	userID, err := s.userService.GetUserIDByName(username)
	if err != nil {
		return "", err
	}
	return security.CreateToken(username, userID, userRole)
}

func (s *jwtAuthService) ValidateToken(token string) (claims jwt.MapClaims, err error) {
	return security.VerifyToken(token)
}
