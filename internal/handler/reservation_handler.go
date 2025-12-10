package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type ReservationHandler struct {
	App *app.App
}

func NewReservationHandler(app *app.App) *ReservationHandler {
	return &ReservationHandler{
		App: app,
	}
}

var ErrUnauthorized = errors.New("Unauthorized")

type CreateReservationRequest struct {
	ShowtimeID uint `json:"showtime_id" binding:"required"`
	SeatID     uint `json:"seat_id" binding:"required"`
}

// @route POST /reservations
func (h *ReservationHandler) CreateReservation(ctx *gin.Context) {
	// Extract user_id from context (guaranteed to exist by RequireAuth middleware)
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		dto.Unauthorized(ctx, "User not authenticated")
		return
	}
	userID := userIDValue.(uint)

	var req CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	err := h.App.ReservationService.Reserve(userID, req.ShowtimeID, req.SeatID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShowtimeNotExist):
			ctx.Error(err)
			dto.NotFound(ctx, "Showtime not found")
		case errors.Is(err, service.ErrNoTicketsAvailable):
			ctx.Error(err)
			dto.Conflict(ctx, "NO_TICKETS", "No tickets available")
		case errors.Is(err, service.ErrAlreadyReserved):
			ctx.Error(err)
			dto.Conflict(ctx, "ALREADY_RESERVED", "You have already reserved this showtime")
		default:
			ctx.Error(err)
			dto.InternalServerError(ctx, "Failed to create reservation")
		}
		return
	}

	dto.SuccessWithMessage(ctx, http.StatusCreated, nil, "Reservation created successfully")
}

// @route GET /reservations/me
func (h *ReservationHandler) GetMyReservations(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		dto.Unauthorized(ctx, "User not authenticated")
		return
	}

	reservations, err := h.App.ReservationService.GetReservationsByUserID(userID)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to retrieve reservations")
		return
	}

	dto.Success(ctx, http.StatusOK, reservations)
}

func getUserIDFromContext(ctx *gin.Context) (uint, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, errors.New("unauthorized")
	}
	return userID.(uint), nil
}

// @route DELETE /reservations/:id
func (h *ReservationHandler) CancelReservation(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		dto.Unauthorized(ctx, "User not authenticated")
		return
	}

	idParam := ctx.Param("id")
	reservationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid reservation ID")
		return
	}

	// Verify that the reservation belongs to the user
	reservation, err := h.App.ReservationService.GetReservationByID(uint(reservationID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.Error(err)
			dto.NotFound(ctx, "Reservation not found")
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to retrieve reservation")
		return
	}

	if reservation.UserID != userID {
		ctx.Error(err)
		dto.Forbidden(ctx, "You are not allowed to cancel this reservation")
		return
	}

	// Cancel the reservation
	err = h.App.ReservationService.CancelReservation(uint(reservationID))
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to cancel reservation")
		return
	}

	dto.SuccessWithMessage(ctx, http.StatusOK, nil, "Reservation cancelled successfully")
}
