package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

type MovieService interface {
	CreateMovie(movie *model.Movie) error
	// UpdateMovie only allows to update the movie having no related showtime
	UpdateMovie(movie *model.Movie) error
	// DeleteMovieByID only allows to delete the movie having no related showtime
	DeleteMovieByID(id uint) error
	GetMovieByID(id uint) (*model.Movie, error)
	GetMovieByTitle(title string) (*model.Movie, error)
	GetAllMovies() ([]model.Movie, error)
}

type movieService struct {
	db                  *gorm.DB
	repo                repository.MovieRepo
	showtimeService     ShowtimeService
	showtimeSeatService ShowtimeSeatService
}

var _ MovieService = (*movieService)(nil)

func NewMovieService(db *gorm.DB, movieRepo repository.MovieRepo, showtimeSeatService ShowtimeSeatService,
	showtimeService ShowtimeService) *movieService {
	return &movieService{
		db:                  db,
		repo:                movieRepo,
		showtimeService:     showtimeService,
		showtimeSeatService: showtimeSeatService,
	}
}

func (s *movieService) CreateMovie(movie *model.Movie) error {
	if err := s.repo.Create(movie); err != nil {
		return err
	}
	return nil
}

var ErrRelatedResourceExists = errors.New("There's are related resources, so can't change")

// Update movie by ID
func (s *movieService) UpdateMovie(movie *model.Movie) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Verify the movie with this ID exists
		existingMovie, err := s.repo.WithTx(tx).GetByID(uint(movie.ID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		// Not allowed to update if related showtime exists
		relatedShowtimes, err := s.showtimeService.GetShowtimesByMovieIDTx(tx, movie.ID)
		if err != nil {
			return err
		}
		if len(relatedShowtimes) != 0 {
			return ErrRelatedResourceExists
		}

		// check if the new title is already used by another
		// because the title needs to be unique
		if existingMovie.Title != movie.Title {
			anotherMovie, err := s.repo.WithTx(tx).GetByTitle(movie.Title)
			if err == nil && anotherMovie != nil && anotherMovie.ID != movie.ID {
				return ErrAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		return s.repo.WithTx(tx).Update(*movie)
	})
}

func (s *movieService) DeleteMovieByID(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Not allowed to delete if related showtime exists
		relatedShowtimes, err := s.showtimeService.GetShowtimesByMovieIDTx(tx, id)
		if err != nil {
			return err
		}
		if len(relatedShowtimes) != 0 {
			return ErrRelatedResourceExists
		}

		return s.repo.WithTx(tx).DeleteByID(id)
	})
}

func (s *movieService) GetMovieByID(id uint) (*model.Movie, error) {
	movie, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *movieService) GetMovieByTitle(title string) (*model.Movie, error) {
	movie, err := s.repo.GetByTitle(title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *movieService) GetAllMovies() ([]model.Movie, error) {
	movies, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return movies, nil
}
