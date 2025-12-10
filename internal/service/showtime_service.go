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
	GetShowtimeByID(showtimeID uint) (*model.Showtime, error)
	GetShowtimesByMovieID(movieID uint) ([]model.Showtime, error)
	GetShowtimesByMovieIDTx(tx *gorm.DB, movieID uint) ([]model.Showtime, error)
	GetShowtimesByHallID(hallID uint) ([]model.Showtime, error)
	GetShowtimesByHallIDTx(tx *gorm.DB, hallID uint) ([]model.Showtime, error)
	GetAllShowtimes() ([]model.Showtime, error)
}

type showtimeService struct {
	db   *gorm.DB
	repo repository.ShowtimeRepo
}

var _ ShowtimeService = (*showtimeService)(nil)

func NewShowtimeService(db *gorm.DB, showtimeRepo repository.ShowtimeRepo) *showtimeService {
	return &showtimeService{
		db:   db,
		repo: showtimeRepo,
	}
}

func (s *showtimeService) CreateShowtime(movieID uint, startTime time.Time, hallID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		showtime := &model.Showtime{
			MovieID: uint(movieID),
			StartAt: startTime,
			HallID:  uint(hallID),
		}
		return s.repo.WithTx(tx).Create(showtime)
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
