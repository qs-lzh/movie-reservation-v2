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

type HallHandler struct {
	App *app.App
}

func NewHallHandler(app *app.App) *HallHandler {
	return &HallHandler{
		App: app,
	}
}

// @route GET /halls
func (h *HallHandler) GetAllHalls(ctx *gin.Context) {
	halls, err := h.App.HallService.GetAllHalls()
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get all halls")
		return
	}
	dto.Success(ctx, http.StatusOK, halls)
}

// @route GET /halls/:id
func (h *HallHandler) GetHallByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid hall id")
		return
	}
	hall, err := h.App.HallService.GetHallByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.Error(err)
			dto.NotFound(ctx, "Hall not exists")
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get hall")
		return
	}
	dto.Success(ctx, http.StatusOK, hall)
}

type CreateHallRequest struct {
	Name      string `json:"name"`
	SeatCount int    `json:"seat_count"`
	Rows      int    `json:"rows"`
	Cols      int    `json:"cols"`
}

// @route POST /halls
func (h *HallHandler) CreateHall(ctx *gin.Context) {
	var req CreateHallRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	hall, err := h.App.HallService.GetHallByName(req.Name)
	if err == nil && hall != nil {
		dto.Conflict(ctx, "HALL_EXISTS", fmt.Sprintf("Hall %s already exists", req.Name))
		return
	}

	hall = &model.Hall{
		Name:      req.Name,
		SeatCount: req.SeatCount,
		Rows:      req.Rows,
		Cols:      req.Cols,
	}

	err = h.App.HallService.CreateHall(hall)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to create hall")
		return
	}

	dto.Success(ctx, http.StatusCreated, hall)
}

type UpdateHallRequest struct {
	Name      string `json:"name"`
	SeatCount int    `json:"seat_count"`
	Rows      int    `json:"rows"`
	Cols      int    `json:"cols"`
}

// @route PUT /halls/:id
func (h *HallHandler) UpdateHall(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid hall id")
		return
	}

	var req UpdateHallRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	// check if the hall exists
	existingHall, err := h.App.HallService.GetHallByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.Error(err)
			dto.NotFound(ctx, "Hall not exists")
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get hall")
		return
	}

	existingHall.Name = req.Name
	existingHall.SeatCount = req.SeatCount
	existingHall.Rows = req.Rows
	existingHall.Cols = req.Cols

	err = h.App.HallService.UpdateHall(existingHall)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to update hall")
		return
	}

	dto.Success(ctx, http.StatusOK, existingHall)
}

// @route DELETE /halls/:id
func (h *HallHandler) DeleteHall(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid hall id")
		return
	}

	err = h.App.HallService.DeleteHallByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.Error(err)
			dto.NotFound(ctx, "Hall not exists")
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to delete hall")
		return
	}

	dto.SuccessWithMessage(ctx, http.StatusOK, nil, "Hall deleted successfully")
}
