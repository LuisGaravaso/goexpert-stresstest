package repository

import (
	"context"
	"stresstest/internal/entity"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the RepositoryInterface
// This is used for testing purposes when Saving Works
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, run *entity.TestRun) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}
