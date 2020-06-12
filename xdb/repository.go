package xdb

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Repository struct {
	DB *gorm.DB
}

func (d *Repository) Create(ctx context.Context, entity interface{}) error {
	return d.CreateTx(ctx, d.DB, entity)
}

func (d *Repository) Save(ctx context.Context, entity interface{}) error {
	return d.SaveTx(ctx, d.DB, entity)
}

func (d *Repository) CreateTx(ctx context.Context, tx *gorm.DB, entity interface{}) error {
	return Trace(ctx, tx).Create(entity).Error
}

func (d *Repository) SaveTx(ctx context.Context, tx *gorm.DB, entity interface{}) error {
	err := Trace(ctx, tx).Save(entity).Error
	return errors.WithStack(err)
}
