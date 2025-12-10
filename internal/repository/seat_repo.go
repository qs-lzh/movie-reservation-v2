package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

type SeatRepo interface {
	WithTx(tx *gorm.DB) SeatRepo
	Create(seat *model.Seat) error
	CreateBatch(seats []model.Seat) error
	GetByID(id uint) (*model.Seat, error)
	GetByHallID(hallID uint) ([]model.Seat, error)
	DeleteByID(id uint) error
}

type seatRepoGorm struct {
	db *gorm.DB
}

var _ SeatRepo = (*seatRepoGorm)(nil)

func NewSeatRepoGorm(db *gorm.DB) *seatRepoGorm {
	return &seatRepoGorm{
		db: db,
	}
}

func (r *seatRepoGorm) WithTx(tx *gorm.DB) SeatRepo {
	return &seatRepoGorm{
		db: tx,
	}
}

func (r *seatRepoGorm) Create(seat *model.Seat) error {
	ctx := context.Background()
	if err := gorm.G[model.Seat](r.db).Create(ctx, seat); err != nil {
		return err
	}
	return nil
}

func (r *seatRepoGorm) CreateBatch(seats []model.Seat) error {
	ctx := context.Background()
	if err := gorm.G[model.Seat](r.db).CreateInBatches(ctx, &seats, len(seats)); err != nil {
		return err
	}
	return nil
}

func (r *seatRepoGorm) GetByID(id uint) (*model.Seat, error) {
	ctx := context.Background()
	seat, err := gorm.G[model.Seat](r.db).Where(&model.Seat{ID: id}).First(ctx)
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *seatRepoGorm) GetByHallID(hallID uint) ([]model.Seat, error) {
	ctx := context.Background()
	seats, err := gorm.G[model.Seat](r.db).Where(&model.Seat{}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *seatRepoGorm) DeleteByID(id uint) error {
	ctx := context.Background()
	_, err := gorm.G[model.Seat](r.db).Where(&model.Seat{ID: id}).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
