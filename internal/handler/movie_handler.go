package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type MovieHandler struct {
	App *app.App
}

func NewMovieHandler(app *app.App) *MovieHandler {
	return &MovieHandler{
		App: app,
	}
}

// @route GET /movies
func (h *MovieHandler) GetAllMovies(ctx *gin.Context) {
	movies, err := h.App.MovieService.GetAllMovies()
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get all movies")
		return
	}
	dto.Success(ctx, http.StatusOK, movies)
}

// @route GET /movies/:id
func (h *MovieHandler) GetMovieByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid movie id")
		return
	}
	movie, err := h.App.MovieService.GetMovieByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.Error(err)
			dto.NotFound(ctx, "Movie not exists")
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get movie")
		return
	}
	dto.Success(ctx, http.StatusOK, movie)
}

func (h *MovieHandler) GetMovieShowtimes(ctx *gin.Context) {
	idParam := ctx.Param("id")
	movieID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid movie id")
		return
	}
	showtimes, err := h.App.ShowtimeService.GetShowtimesByMovieID(uint(movieID))
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get showtimes")
		return
	}
	dto.Success(ctx, http.StatusOK, showtimes)
}

type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// @route POST /movies
func (h *MovieHandler) CreateMovie(ctx *gin.Context) {
	var req CreateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	// confirm no movie is the same title
	anotherMovie, err := h.App.MovieService.GetMovieByTitle(req.Title)
	if err == nil && anotherMovie != nil {
		dto.Conflict(ctx, "Movie_EXISTS", fmt.Sprintf("Movie %s already exists", req.Title))
		return
	}

	movie := &model.Movie{
		Title:       req.Title,
		Description: req.Description,
	}

	err = h.App.MovieService.CreateMovie(movie)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to create movie")
		return
	}

	dto.Success(ctx, http.StatusCreated, movie)
}

type UpdateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// @route PUT /movies/:id
func (h *MovieHandler) UpdateMovie(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid movie id")
		return
	}

	var req UpdateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	// check if the movie exists
	existingMovie, err := h.App.MovieService.GetMovieByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.Error(err)
			dto.NotFound(ctx, "Movie not exists")
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get movie")
		return
	}

	existingMovie.Title = req.Title
	existingMovie.Description = req.Description

	err = h.App.MovieService.UpdateMovie(existingMovie)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to update movie")
		return
	}

	dto.Success(ctx, http.StatusOK, existingMovie)
}

// @route DELETE /movies/:id
func (h *MovieHandler) DeleteMovie(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid movie id")
		return
	}

	err = h.App.MovieService.DeleteMovieByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.Error(err)
			dto.NotFound(ctx, "Movie not exists")
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to delete movie")
		return
	}

	dto.SuccessWithMessage(ctx, http.StatusOK, nil, "Movie deleted successfully")
}
