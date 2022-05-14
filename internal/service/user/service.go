package user

import (
	"context"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/bipen2001/go-user-assignment-go/internal/service/user/model"
)

type service struct {
	repo model.Repository
}

func NewService(repo model.Repository) model.Service {
	return service{repo}
}

func (s service) Get(ctx context.Context) ([]entity.User, error) {
	users, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s service) GetById(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	usr, err := s.repo.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (s service) Update(ctx context.Context, user entity.User) (*entity.User, error) {
	usr, err := s.repo.Update(ctx, user)

	if err != nil {
		return nil, err
	}

	return usr, nil

}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil

}
