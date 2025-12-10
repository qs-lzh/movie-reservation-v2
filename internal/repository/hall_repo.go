package repository

import (
	"context"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"gorm.io/gorm"
)

type HallRepo interface {
	WithTx(tx *gorm.DB) HallRepo
	Create(hall *model.Hall) error
	GetByID(id uint) (*model.Hall, error)
	GetByName(name string) (*model.Hall, error)
	DeleteByID(id uint) error
	ListAll() ([]model.Hall, error)
	Update(*model.Hall) error
}

type hallRepoGorm struct {
	db *gorm.DB
}

var _ HallRepo = (*hallRepoGorm)(nil)

func (r *hallRepoGorm) WithTx(tx *gorm.DB) HallRepo {
	return &hallRepoGorm{
		db: tx,
	}
}

func NewHallRepoGorm(db *gorm.DB) *hallRepoGorm {
	return &hallRepoGorm{
		db: db,
	}
}

func (r *hallRepoGorm) Create(hall *model.Hall) error {
	ctx := context.Background()
	context.Background()
	if err := gorm.G[model.Hall](r.db).Create(ctx, hall); err != nil {
		return err
	}
	return nil
}

func (r *hallRepoGorm) GetByID(id uint) (*model.Hall, error) {
	ctx := context.Background()
	hall, err := gorm.G[model.Hall](r.db).Where(&model.Hall{ID: id}).First(ctx)
	if err != nil {
		return nil, err
	}
	return &hall, nil
}

func (r *hallRepoGorm) GetByName(name string) (*model.Hall, error) {
	ctx := context.Background()
	hall, err := gorm.G[model.Hall](r.db).Where(&model.Hall{Name: name}).First(ctx)
	if err != nil {
		return nil, err
	}
	return &hall, nil
}

func (r *hallRepoGorm) DeleteByID(id uint) error {
	ctx := context.Background()
	_, err := gorm.G[model.Hall](r.db).Where(&model.Hall{ID: id}).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *hallRepoGorm) ListAll() ([]model.Hall, error) {
	ctx := context.Background()
	halls, err := gorm.G[model.Hall](r.db).Find(ctx)
	if err != nil {
		return nil, err
	}
	return halls, nil
}

// before use Update, please confirm the existance of the hall
func (r *hallRepoGorm) Update(hall *model.Hall) error {
	ctx := context.Background()
	if _, err := gorm.G[model.Hall](r.db).Updates(ctx, *hall); err != nil {
		return err
	}
	return nil
}
