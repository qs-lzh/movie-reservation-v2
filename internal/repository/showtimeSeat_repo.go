package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

type ShowtimeSeatRepo interface {
	WithTx(tx *gorm.DB) ShowtimeSeatRepo
	Create(showtimeSeat *model.ShowtimeSeat) error
	CreateBatch(showtimeSeats []model.ShowtimeSeat) error
	GetByID(id uint) (*model.ShowtimeSeat, error)
	GetByShowtimeID(showtimeID uint) ([]model.ShowtimeSeat, error)
	GetBySeatID(seatID uint) ([]model.ShowtimeSeat, error)
	GetByShowIDSeatID(showtimeID, seatID uint) (*model.ShowtimeSeat, error)
	GetByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error)
	Update(id uint, showtimeSeat *model.ShowtimeSeat) error
	DeleteByID(id uint) error
}

type showtimeSeatRepoGorm struct {
	db *gorm.DB
}

var _ ShowtimeSeatRepo = (*showtimeSeatRepoGorm)(nil)

func NewShowtimeSeatRepoGorm(db *gorm.DB) *showtimeSeatRepoGorm {
	return &showtimeSeatRepoGorm{
		db: db,
	}
}

func (r *showtimeSeatRepoGorm) WithTx(tx *gorm.DB) ShowtimeSeatRepo {
	return &showtimeSeatRepoGorm{
		db: tx,
	}
}

func (r *showtimeSeatRepoGorm) Create(showtimeSeat *model.ShowtimeSeat) error {
	ctx := context.Background()
	if err := gorm.G[model.ShowtimeSeat](r.db).Create(ctx, showtimeSeat); err != nil {
		return err
	}
	return nil
}

func (r *showtimeSeatRepoGorm) CreateBatch(showtimeSeats []model.ShowtimeSeat) error {
	ctx := context.Background()
	if err := gorm.G[model.ShowtimeSeat](r.db).CreateInBatches(ctx, &showtimeSeats, len(showtimeSeats)); err != nil {
		return err
	}
	return nil
}

func (r *showtimeSeatRepoGorm) GetByID(id uint) (*model.ShowtimeSeat, error) {
	ctx := context.Background()
	showtimeSeat, err := gorm.G[model.ShowtimeSeat](r.db).Where(&model.ShowtimeSeat{ID: id}).First(ctx)
	if err != nil {
		return nil, err
	}
	return &showtimeSeat, nil
}

func (r *showtimeSeatRepoGorm) GetByShowtimeID(showtimeID uint) ([]model.ShowtimeSeat, error) {
	ctx := context.Background()
	showtimeSeats, err := gorm.G[model.ShowtimeSeat](r.db).Where(&model.ShowtimeSeat{ShowtimeID: showtimeID}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return showtimeSeats, nil
}

func (r *showtimeSeatRepoGorm) GetBySeatID(seatID uint) ([]model.ShowtimeSeat, error) {
	ctx := context.Background()
	showtimeSeats, err := gorm.G[model.ShowtimeSeat](r.db).Where(&model.ShowtimeSeat{SeatID: seatID}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return showtimeSeats, nil
}

func (r *showtimeSeatRepoGorm) GetByShowIDSeatID(showtimeID, seatID uint) (*model.ShowtimeSeat, error) {
	ctx := context.Background()
	showtimeSeat, err := gorm.G[model.ShowtimeSeat](r.db).Where(&model.ShowtimeSeat{ShowtimeID: showtimeID, SeatID: seatID}).First(ctx)
	if err != nil {
		return nil, err
	}
	return &showtimeSeat, nil
}

func (r *showtimeSeatRepoGorm) GetByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error) {
	ctx := context.Background()
	showtimeSeats, err := gorm.G[model.ShowtimeSeat](r.db).Where(&model.ShowtimeSeat{Status: status}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return showtimeSeats, nil
}

func (r *showtimeSeatRepoGorm) Update(id uint, showtimeSeat *model.ShowtimeSeat) error {
	ctx := context.Background()
	if _, err := gorm.G[model.ShowtimeSeat](r.db).Updates(ctx, model.ShowtimeSeat{ID: id}); err != nil {
		return err
	}
	return nil
}

func (r *showtimeSeatRepoGorm) DeleteByID(id uint) error {
	ctx := context.Background()
	if _, err := gorm.G[model.ShowtimeSeat](r.db).Where(&model.ShowtimeSeat{ID: id}).Delete(ctx); err != nil {
		return err
	}
	return nil
}
