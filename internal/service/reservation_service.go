package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type ReservationService interface {
	Reserve(userID, showtimeID, seatID uint) error
	CancelReservation(reservationID uint) error
	GetRemainingTicketsTx(tx *gorm.DB, showtime *model.Showtime) (int, error)
	GetReservationsByUserID(userID uint) ([]model.Reservation, error)
	GetReservationsByUserIDTx(tx *gorm.DB, userID uint) ([]model.Reservation, error)
	GetReservationByID(reservationID uint) (*model.Reservation, error)
}

type reservationService struct {
	db                  *gorm.DB
	repo                repository.ReservationRepo
	showtimeRepo        repository.ShowtimeRepo
	hallRepo            repository.HallRepo
	showtimeSeatService ShowtimeSeatService
}

var _ ReservationService = (*reservationService)(nil)

func NewReservationService(db *gorm.DB, reservationRepo repository.ReservationRepo,
	showtimeRepo repository.ShowtimeRepo, hallRepo repository.HallRepo, showtimeSeatService ShowtimeSeatService) *reservationService {
	return &reservationService{
		db:                  db,
		repo:                reservationRepo,
		showtimeRepo:        showtimeRepo,
		hallRepo:            hallRepo,
		showtimeSeatService: showtimeSeatService,
	}
}

func (s *reservationService) Reserve(userID, showtimeID, seatID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// check if showtime exists
		showtime, err := s.showtimeRepo.WithTx(tx).GetByID(showtimeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrShowtimeNotExist
			}
			return err
		}

		// check if there's tickets available
		if _, err := s.GetRemainingTicketsTx(tx, showtime); err != nil {
			return err
		}

		// check if the user already have the same reservation
		reservations, err := s.repo.WithTx(tx).GetByUserID(userID)
		if err != nil {
			return err
		}
		for _, reservation := range reservations {
			if reservation.ShowtimeID == showtimeID {
				return ErrAlreadyReserved
			}
		}

		// reserve
		if err = s.repo.WithTx(tx).Create(&model.Reservation{
			ShowtimeID: showtimeID,
			SeatID:     seatID,
			UserID:     userID,
		}); err != nil {
			return err
		}

		// change showtimeSeat status
		showtimeSeat, err := s.showtimeSeatService.GetShowtimeSeatByShowtimeIDSeatIDTx(tx, showtimeID, seatID)
		if err != nil {
			return err
		}
		return s.showtimeSeatService.UpdateShowtimeSeatStatusToLockedTx(tx, showtimeSeat.ID)
	})
}

func (s *reservationService) CancelReservation(reservationID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		reservation, err := s.repo.WithTx(tx).GetByID(reservationID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if err := s.repo.WithTx(tx).DeleteByID(reservationID); err != nil {
			return err
		}

		// change showtimeSeat status
		showtimeSeat, err := s.showtimeSeatService.GetShowtimeSeatByShowtimeIDSeatIDTx(tx, reservation.ShowtimeID, reservation.SeatID)
		if err != nil {
			return err
		}
		return s.showtimeSeatService.UpdateShowtimeSeatStatusToAvailableTx(tx, showtimeSeat.ID)
	})
}

func (s *reservationService) GetRemainingTicketsTx(tx *gorm.DB, showtime *model.Showtime) (int, error) {
	var remainingTickets int
	err := s.db.Transaction(func(tx *gorm.DB) error {
		reservations, err := s.repo.WithTx(tx).GetByShowtimeID(showtime.ID)
		if err != nil {
			return err
		}
		hall, err := s.hallRepo.WithTx(tx).GetByID(showtime.HallID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		remainingTickets = hall.SeatCount - len(reservations)
		if remainingTickets <= 0 {
			return ErrNoTicketsAvailable
		}
		return nil
	})

	return remainingTickets, err
}

func (s *reservationService) GetReservationsByUserID(userID uint) ([]model.Reservation, error) {
	return s.GetReservationsByUserIDTx(s.db, userID)
}

func (s *reservationService) GetReservationsByUserIDTx(tx *gorm.DB, userID uint) ([]model.Reservation, error) {
	return s.repo.WithTx(tx).GetByUserID(userID)
}
func (s *reservationService) GetReservationByID(reservationID uint) (*model.Reservation, error) {
	reservation, err := s.repo.GetByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return reservation, nil
}
