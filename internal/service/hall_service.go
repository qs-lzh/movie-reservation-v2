package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

type HallService interface {
	CreateHall(hall *model.Hall) error
	UpdateHall(hall *model.Hall) error
	DeleteHallByID(id uint) error
	GetHallByID(id uint) (*model.Hall, error)
	GetHallByName(name string) (*model.Hall, error)
	GetAllHalls() ([]model.Hall, error)
}

type hallService struct {
	db              *gorm.DB
	repo            repository.HallRepo
	seatService     SeatService
	showtimeService ShowtimeService
}

var _ HallService = (*hallService)(nil)

func NewHallService(db *gorm.DB, hallRepo repository.HallRepo, seatService SeatService,
	showtimeService ShowtimeService) *hallService {
	return &hallService{
		db:              db,
		repo:            hallRepo,
		seatService:     seatService,
		showtimeService: showtimeService,
	}
}

func (s *hallService) CreateHall(hall *model.Hall) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.WithTx(tx).Create(hall); err != nil {
			return err
		}
		return s.seatService.InitSeatsForHallTx(tx, hall)
	})
}

func (s *hallService) UpdateHall(hall *model.Hall) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// verify no related Showtime
		relatedShowtimes, err := s.showtimeService.GetShowtimesByHallIDTx(tx, hall.ID)
		if err != nil {
			return err
		}
		if len(relatedShowtimes) != 0 {
			return ErrRelatedResourceExists
		}

		// verify that the hall with this ID exists
		existinghall, err := s.repo.WithTx(tx).GetByID(uint(hall.ID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		// check if the new title is already used by another
		// because the title needs to be unique
		if existinghall.Name != hall.Name {
			anotherhall, err := s.repo.WithTx(tx).GetByName(hall.Name)
			if err == nil && anotherhall != nil && anotherhall.ID != hall.ID {
				return ErrAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		return s.repo.WithTx(tx).Update(hall)
	})
}

func (s *hallService) DeleteHallByID(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// verify no related showtime exists
		relatedShowtimes, err := s.showtimeService.GetShowtimesByHallIDTx(tx, id)
		if err != nil {
			return err
		}
		if len(relatedShowtimes) != 0 {
			return ErrRelatedResourceExists
		}

		return s.repo.WithTx(tx).DeleteByID(id)
	})
}

func (s *hallService) GetHallByID(id uint) (*model.Hall, error) {
	hall, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return hall, nil
}

func (s *hallService) GetHallByName(name string) (*model.Hall, error) {
	hall, err := s.repo.GetByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return hall, nil
}

func (s *hallService) GetAllHalls() ([]model.Hall, error) {
	halls, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return halls, nil
}
