package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"github.com/qs-lzh/movie-reservation/internal/security"
)

type UserService interface {
	CreateUser(userName, password string, role model.UserRole) error
	DeleteUser(userName string, password string) error
	ValidateUser(userName string, password string) (bool, error)
	GetUserRoleByName(userName string) (model.UserRole, error)
	GetUserIDByName(userName string) (uint, error)
}

type userService struct {
	db                 *gorm.DB
	hasher             security.PasswordHasher
	repo               repository.UserRepo
	reservationService ReservationService
}

var _ UserService = (*userService)(nil)

func NewUserService(db *gorm.DB, userRepo repository.UserRepo, reservationService ReservationService) *userService {
	return &userService{
		db:                 db,
		hasher:             security.NewBcryptHasher(10),
		repo:               userRepo,
		reservationService: reservationService,
	}
}

func (s *userService) CreateUser(userName, password string, role model.UserRole) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		_, err := s.repo.WithTx(tx).GetByName(userName)
		if err == nil {
			return ErrAlreadyExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		hash, err := s.hasher.Hash(password)
		if err != nil {
			return err
		}
		return s.repo.WithTx(tx).Create(&model.User{
			Name:           userName,
			HashedPassword: hash,
			Role:           role,
		})
	})
}

func (s *userService) DeleteUser(userName string, password string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.WithTx(tx).GetByName(userName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if err = s.hasher.Compare(user.HashedPassword, password); err != nil {
			return ErrInvalidCredential
		}

		// ensure no related reservation
		s.reservationService.GetReservationsByUserIDTx(tx, user.ID)

		return s.repo.WithTx(tx).DeleteByName(userName)
	})
}

func (s *userService) ValidateUser(userName string, password string) (bool, error) {
	user, err := s.repo.GetByName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if err = s.hasher.Compare(user.HashedPassword, password); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *userService) GetUserRoleByName(userName string) (model.UserRole, error) {
	user, err := s.repo.GetByName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return user.Role, nil
}

func (s *userService) GetUserIDByName(userName string) (uint, error) {
	user, err := s.repo.GetByName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return uint(user.ID), nil
}
