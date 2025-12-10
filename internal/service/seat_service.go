package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

/*
* SeatService should only be used by hall service
* So SeatService do not examine if the operand is allowed or is safe
 */

type SeatService interface {
	CreateSeat(seat *model.Seat) error
	InitSeatsForHall(hall *model.Hall) error
	InitSeatsForHallTx(tx *gorm.DB, hall *model.Hall) error
	GetSeatByID(id uint) (*model.Seat, error)
	GetSeatsByHallID(hallID uint) ([]model.Seat, error)
	GetSeatsByHallIDTx(tx *gorm.DB, hallID uint) ([]model.Seat, error)
	// DeleteSeatByID do not examine if it is allowed to delete,
	DeleteSeatByID(id uint) error
}

type seatService struct {
	db   *gorm.DB
	repo repository.SeatRepo
}

func NewseatService(db *gorm.DB, seatRepo repository.SeatRepo) *seatService {
	return &seatService{
		db:   db,
		repo: seatRepo,
	}
}

var _ SeatService = (*seatService)(nil)

func (s *seatService) CreateSeat(seat *model.Seat) error {
	return s.repo.Create(seat)
}

func (s *seatService) InitSeatsForHall(hall *model.Hall) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.InitSeatsForHallTx(s.db, hall)
	})
}

func (s *seatService) InitSeatsForHallTx(tx *gorm.DB, hall *model.Hall) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		hallID := hall.ID
		rows, cols := hall.Rows, hall.Cols
		var seats []model.Seat
		for row := range rows {
			for col := range cols {
				seats = append(seats, model.Seat{
					HallID: hallID,
					Row:    row,
					Col:    col,
				})
			}
		}

		return s.repo.WithTx(tx).CreateBatch(seats)
	})
}

func (s *seatService) GetSeatByID(id uint) (*model.Seat, error) {
	seat, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return seat, nil
}

func (s *seatService) GetSeatsByHallID(hallID uint) ([]model.Seat, error) {
	return s.GetSeatsByHallIDTx(s.db, hallID)
}

func (s *seatService) GetSeatsByHallIDTx(tx *gorm.DB, hallID uint) ([]model.Seat, error) {
	return s.repo.WithTx(tx).GetByHallID(hallID)
}

func (s *seatService) DeleteSeatByID(id uint) error {
	if err := s.repo.DeleteByID(id); err != nil {
		return err
	}
	return nil
}
