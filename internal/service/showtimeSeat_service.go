package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

/*
* ShowtimeSeatService do not examine if the operand is allowed or is safe
 */

type ShowtimeSeatService interface {
	CreateShowtimeSeat(showtimeSeat *model.ShowtimeSeat) error
	InitShowtimeSeatsForShowtimeTx(tx *gorm.DB, showtime *model.Showtime) error
	GetShowtimeSeatByID(id uint) (*model.ShowtimeSeat, error)
	GetShowtimeSeatByShowtimeIDSeatID(showtimeID, seatID uint) (*model.ShowtimeSeat, error)
	GetShowtimeSeatByShowtimeIDSeatIDTx(tx *gorm.DB, showtimeID, seatID uint) (*model.ShowtimeSeat, error)
	GetShowtimeSeatsByShowtimeID(showtimeID uint) ([]model.ShowtimeSeat, error)
	GetShowtimeSeatsByShowtimeIDTx(tx *gorm.DB, showtimeID uint) ([]model.ShowtimeSeat, error)
	GetShowtimeSeatsByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error)
	updateShowtimeSeatStatusTx(tx *gorm.DB, id uint, targetStatus model.ShowtimeSeatStatus) error
	UpdateShowtimeSeatStatusToAvailableTx(tx *gorm.DB, id uint) error
	UpdateShowtimeSeatStatusToLockedTx(tx *gorm.DB, id uint) error
	UpdateShowtimeSeatStatusToSoldTx(tx *gorm.DB, id uint) error
	DeleteShowtimeSeatByID(id uint) error
}

type showtimeSeatService struct {
	db          *gorm.DB
	repo        repository.ShowtimeSeatRepo
	seatService SeatService
}

func NewShowtimeSeatService(db *gorm.DB, showtimeSeatRepo repository.ShowtimeSeatRepo, seatService SeatService) *showtimeSeatService {
	return &showtimeSeatService{
		db:          db,
		repo:        showtimeSeatRepo,
		seatService: seatService,
	}
}

var _ ShowtimeSeatService = (*showtimeSeatService)(nil)

func (s *showtimeSeatService) CreateShowtimeSeat(showtimeSeat *model.ShowtimeSeat) error {
	if err := s.repo.Create(showtimeSeat); err != nil {
		return err
	}
	return nil
}

func (s *showtimeSeatService) InitShowtimeSeatsForShowtimeTx(tx *gorm.DB, showtime *model.Showtime) error {
	hallID := showtime.HallID
	seats, err := s.seatService.GetSeatsByHallIDTx(tx, hallID)
	if err != nil {
		return err
	}
	var showtimeSeats []model.ShowtimeSeat
	for _, seat := range seats {
		showtimeSeats = append(showtimeSeats, model.ShowtimeSeat{
			ShowtimeID: showtime.ID,
			SeatID:     seat.ID,
			Status:     model.StatusAvailable,
		})
	}
	return s.repo.WithTx(tx).CreateBatch(showtimeSeats)
}

func (s *showtimeSeatService) GetShowtimeSeatByID(id uint) (*model.ShowtimeSeat, error) {
	showtimeSeat, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return showtimeSeat, nil
}

func (s *showtimeSeatService) GetShowtimeSeatByShowtimeIDSeatID(showtimeID, seatID uint) (*model.ShowtimeSeat, error) {
	return s.GetShowtimeSeatByShowtimeIDSeatIDTx(s.db, showtimeID, seatID)
}

func (s *showtimeSeatService) GetShowtimeSeatByShowtimeIDSeatIDTx(tx *gorm.DB, showtimeID, seatID uint) (*model.ShowtimeSeat, error) {
	return s.repo.WithTx(tx).GetByShowIDSeatID(showtimeID, seatID)
}

func (s *showtimeSeatService) GetShowtimeSeatsByShowtimeID(showtimeID uint) ([]model.ShowtimeSeat, error) {
	return s.GetShowtimeSeatsByShowtimeIDTx(s.db, showtimeID)
}

func (s *showtimeSeatService) GetShowtimeSeatsByShowtimeIDTx(tx *gorm.DB, showtimeID uint) ([]model.ShowtimeSeat, error) {
	return s.repo.WithTx(tx).GetByShowtimeID(showtimeID)
}

func (s *showtimeSeatService) GetShowtimeSeatsByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error) {
	return s.repo.GetByStatus(status)
}

var ErrShowtimeSeatNotExist = errors.New("The ShowtimeSeat does not exist")
var ErrShowtimeSeatStatusNotChange = errors.New("The target status of showtimeSeat is the same as the origin")

func (s *showtimeSeatService) updateShowtimeSeatStatusTx(tx *gorm.DB, id uint, targetStatus model.ShowtimeSeatStatus) error {
	existingShowtimeSeat, err := s.repo.WithTx(tx).GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrShowtimeSeatNotExist
		}
		return err
	}

	if targetStatus == existingShowtimeSeat.Status {
		return ErrShowtimeSeatStatusNotChange
	}

	existingShowtimeSeat.Status = targetStatus
	return s.repo.WithTx(tx).Update(id, existingShowtimeSeat)
}
func (s *showtimeSeatService) UpdateShowtimeSeatStatusToAvailableTx(tx *gorm.DB, id uint) error {
	return s.updateShowtimeSeatStatusTx(tx, id, model.StatusAvailable)
}
func (s *showtimeSeatService) UpdateShowtimeSeatStatusToLockedTx(tx *gorm.DB, id uint) error {
	return s.updateShowtimeSeatStatusTx(tx, id, model.StatusLocked)
}
func (s *showtimeSeatService) UpdateShowtimeSeatStatusToSoldTx(tx *gorm.DB, id uint) error {
	return s.updateShowtimeSeatStatusTx(tx, id, model.StatusSold)
}

func (s *showtimeSeatService) DeleteShowtimeSeatByID(id uint) error {
	return s.repo.DeleteByID(id)
}
