package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

type MovieService interface {
	CreateMovie(movie *model.Movie) error
	GetMovieByID(id uint) (*model.Movie, error)
	GetMovieByTitle(title string) (*model.Movie, error)
	GetAllMovies() ([]model.Movie, error)
}

type movieService struct {
	db              *gorm.DB
	repo            repository.MovieRepo
	showtimeService ShowtimeService
}

var _ MovieService = (*movieService)(nil)

func NewMovieService(db *gorm.DB, movieRepo repository.MovieRepo, showtimeService ShowtimeService) *movieService {
	return &movieService{
		db:              db,
		repo:            movieRepo,
		showtimeService: showtimeService,
	}
}

func (s *movieService) CreateMovie(movie *model.Movie) error {
	if err := s.repo.Create(movie); err != nil {
		return err
	}
	return nil
}

var ErrRelatedResourceExists = errors.New("There's are related resources, so can't change")

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
