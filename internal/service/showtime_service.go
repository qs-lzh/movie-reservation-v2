package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type ShowtimeService interface {
	CreateShowtime(movieID uint, startTime time.Time, hallID uint) error
	UpdateShowtime(showtimeID uint, startTime time.Time, hallID uint) error
	DeleteShowtimeByID(showtimeID uint) error
	GetShowtimeByID(showtimeID uint) (*model.Showtime, error)
	GetShowtimesByMovieID(movieID uint) ([]model.Showtime, error)
	GetShowtimesByMovieIDTx(tx *gorm.DB, movieID uint) ([]model.Showtime, error)
	GetShowtimesByHallID(hallID uint) ([]model.Showtime, error)
	GetShowtimesByHallIDTx(tx *gorm.DB, hallID uint) ([]model.Showtime, error)
	GetAllShowtimes() ([]model.Showtime, error)
}

type showtimeService struct {
	db                  *gorm.DB
	repo                repository.ShowtimeRepo
	showtimeSeatService ShowtimeSeatService
}

var _ ShowtimeService = (*showtimeService)(nil)

func NewShowtimeService(db *gorm.DB, showtimeRepo repository.ShowtimeRepo, showtimeSeatService ShowtimeSeatService) *showtimeService {
	return &showtimeService{
		db:                  db,
		repo:                showtimeRepo,
		showtimeSeatService: showtimeSeatService,
	}
}

func (s *showtimeService) CreateShowtime(movieID uint, startTime time.Time, hallID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		showtime := &model.Showtime{
			MovieID: uint(movieID),
			StartAt: startTime,
			HallID:  uint(hallID),
		}
		if err := s.repo.WithTx(tx).Create(showtime); err != nil {
			return err
		}

		return s.showtimeSeatService.InitShowtimeSeatsForShowtimeTx(tx, showtime)
	})
}

func (s *showtimeService) UpdateShowtime(showtimeID uint, startTime time.Time, hallID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Ensure no related ShowtimeSeat
		relatedShowtimeSeats, err := s.showtimeSeatService.GetShowtimeSeatsByShowtimeIDTx(tx, showtimeID)
		if err != nil {
			return err
		}
		if len(relatedShowtimeSeats) != 0 {
			return ErrRelatedResourceExists
		}

		showtime, err := s.repo.WithTx(tx).GetByID(uint(showtimeID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if err := s.repo.WithTx(tx).DeleteByID(uint(showtimeID)); err != nil {
			return err
		}
		showtime.ID = uint(showtimeID)
		showtime.StartAt = startTime
		showtime.HallID = uint(hallID)
		if err := s.repo.WithTx(tx).Create(showtime); err != nil {
			return err
		}
		return nil
	})
}

func (s *showtimeService) DeleteShowtimeByID(showtimeID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Ensure no related ShowtimeSeat
		relatedShowtimeSeats, err := s.showtimeSeatService.GetShowtimeSeatsByShowtimeIDTx(tx, showtimeID)
		if err != nil {
			return err
		}
		if len(relatedShowtimeSeats) != 0 {
			return ErrRelatedResourceExists
		}

		return s.repo.WithTx(tx).DeleteByID(uint(showtimeID))
	})
}

func (s *showtimeService) GetShowtimeByID(showtimeID uint) (*model.Showtime, error) {
	showtime, err := s.repo.GetByID(uint(showtimeID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return showtime, nil
}

func (s *showtimeService) GetShowtimesByMovieID(movieID uint) ([]model.Showtime, error) {
	return s.GetShowtimesByMovieIDTx(s.db, movieID)
}
func (s *showtimeService) GetShowtimesByMovieIDTx(tx *gorm.DB, movieID uint) ([]model.Showtime, error) {
	return s.repo.WithTx(tx).GetByMovieID(movieID)
}

func (s *showtimeService) GetShowtimesByHallID(hallID uint) ([]model.Showtime, error) {
	return s.GetShowtimesByHallIDTx(s.db, hallID)
}
func (s *showtimeService) GetShowtimesByHallIDTx(tx *gorm.DB, hallID uint) ([]model.Showtime, error) {
	return s.repo.WithTx(tx).GetByHallID(hallID)
}

func (s *showtimeService) GetAllShowtimes() ([]model.Showtime, error) {
	return s.repo.ListAll()
}
