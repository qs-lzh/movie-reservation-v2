package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

type ReservationRepo interface {
	WithTx(tx *gorm.DB) ReservationRepo
	Create(reservation *model.Reservation) error
	GetByID(id uint) (*model.Reservation, error)
	DeleteByID(id uint) error
	GetByUserID(userID uint) ([]model.Reservation, error)
	GetByShowtimeID(showtimeID uint) ([]model.Reservation, error)
}

type reservationRepoGorm struct {
	db *gorm.DB
}

var _ ReservationRepo = (*reservationRepoGorm)(nil)

func NewReservationRepoGorm(db *gorm.DB) *reservationRepoGorm {
	return &reservationRepoGorm{
		db: db,
	}
}

func (r *reservationRepoGorm) WithTx(tx *gorm.DB) ReservationRepo {
	return &reservationRepoGorm{
		db: tx,
	}
}

func (r *reservationRepoGorm) Create(reservation *model.Reservation) error {
	ctx := context.Background()
	if err := gorm.G[model.Reservation](r.db).Create(ctx, reservation); err != nil {
		return err
	}
	return nil
}

func (r *reservationRepoGorm) GetByID(id uint) (*model.Reservation, error) {
	ctx := context.Background()
	reservation, err := gorm.G[model.Reservation](r.db).Where(&model.Reservation{ID: id}).First(ctx)
	if err != nil {
		return &model.Reservation{}, err
	}
	return &reservation, nil
}

func (r *reservationRepoGorm) DeleteByID(id uint) error {
	ctx := context.Background()
	_, err := gorm.G[model.Reservation](r.db).Where(&model.Reservation{ID: id}).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *reservationRepoGorm) GetByUserID(userID uint) ([]model.Reservation, error) {
	ctx := context.Background()
	reservations, err := gorm.G[model.Reservation](r.db).Where(&model.Reservation{UserID: userID}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepoGorm) GetByShowtimeID(showtimeID uint) ([]model.Reservation, error) {
	ctx := context.Background()
	reservations, err := gorm.G[model.Reservation](r.db).Where(&model.Reservation{ShowtimeID: showtimeID}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}
