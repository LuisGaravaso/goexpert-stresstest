package repository

import (
	"context"
	"stresstest/internal/entity"
)

// RepositoryInterface is an interface for the repository layer
// This project will not implement it, but it is here to show how it could be done
// This will allow future implementations to be easily swapped

type RepositoryInterface interface {
	Save(ctx context.Context, testRun *entity.TestRun) error
}

type Repository struct{}

func NewRepository() Repository {
	return Repository{}
}

func (r *Repository) Save(ctx context.Context, testRun *entity.TestRun) error {
	return nil
}
