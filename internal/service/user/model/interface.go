package model

import (
	"context"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
)

type Repository interface {
	common
}

type Service interface {
	common
}

type common interface {
	GetById(context.Context, string) (*entity.User, error)
	Get(context.Context) ([]entity.User, error)
	Create(context.Context, entity.User) (*entity.User, error)
	Delete(context.Context, string) error
	Update(context.Context, entity.User) (*entity.User, error)
}
