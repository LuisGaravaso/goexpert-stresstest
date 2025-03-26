package run_test

import (
	"context"
	"errors"
	"stresstest/internal/entity"
	"stresstest/internal/usecase/run"
	"stresstest/mocks/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_RunUseCase_MustFailForInvalidParams(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	uc := run.NewRunUseCase(repo)
	input := run.RunInputDTO{Url: "invalid-url", Requests: 10, Concurrency: 10}
	ctx := context.Background()

	// Act
	_, err := uc.Run(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Error(t, err, entity.ErrInvalidURL)
}

func Test_RunUseCase_MustFailForNegativeRequests(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	uc := run.NewRunUseCase(repo)
	input := run.RunInputDTO{Url: "http://example.com", Requests: -1, Concurrency: 10}
	ctx := context.Background()

	// Act
	_, err := uc.Run(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Error(t, err, entity.ErrNonNegativeRequests)
}

func Test_RunUseCase_MustFailForNegativeConcurrency(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	uc := run.NewRunUseCase(repo)
	input := run.RunInputDTO{Url: "http://example.com", Requests: 10, Concurrency: -1}
	ctx := context.Background()

	// Act
	_, err := uc.Run(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Error(t, err, entity.ErrNonNegativeConcurrency)
}

func Test_RunUseCase_MustPassForValidParams(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	uc := run.NewRunUseCase(repo)
	input := run.RunInputDTO{Url: "http://example.com", Requests: 10, Concurrency: 10}
	ctx := context.Background()

	// Act
	output, err := uc.Run(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, input.Url, output.Url)
	assert.Equal(t, input.Requests, output.Requests)
	assert.Equal(t, input.Concurrency, output.Concurrency)
	assert.NotEmpty(t, output.Id)
	assert.NotEmpty(t, output.TimestampStart)
	assert.NotEmpty(t, output.TimestampEnd)
	assert.GreaterOrEqual(t, output.TestDurationInSeconds, 0)
}

func Test_RunUseCase_MustFailSavingToRepository(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("failed to save")).Once()

	uc := run.NewRunUseCase(repo)
	input := run.RunInputDTO{Url: "https://github.com/LuisGaravaso/goexpert-auction", Requests: 10, Concurrency: 10}
	ctx := context.Background()

	// Act
	_, err := uc.Run(ctx, input)

	// Assert
	assert.Error(t, err)
}

func Test_MustMakeFiveRequestsAndShowData(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	uc := run.NewRunUseCase(repo)
	input := run.RunInputDTO{Url: "https://github.com/LuisGaravaso/goexpert-auction", Requests: 5, Concurrency: 1, ShowData: true}
	ctx := context.Background()

	// Act
	output, err := uc.Run(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, input.Url, output.Url)
	assert.Equal(t, input.Requests, output.Requests)
	assert.Equal(t, input.Concurrency, output.Concurrency)
	assert.NotEmpty(t, output.Id)
	assert.NotEmpty(t, output.TimestampStart)
	assert.NotEmpty(t, output.TimestampEnd)
	assert.GreaterOrEqual(t, output.TestDurationInSeconds, 0)
	assert.NotEmpty(t, output.Data)
	assert.Len(t, output.Data, 5)
	assert.NotEmpty(t, output.Data[0].DurationInMs)
	assert.NotEmpty(t, output.Data[0].StatusCode)
	assert.NotEmpty(t, output.Data[4].DurationInMs)
	assert.NotEmpty(t, output.Data[4].StatusCode)
	t.Log(output)
}

func Test_MustMakeFiveRequestsAndSupressData(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	uc := run.NewRunUseCase(repo)
	input := run.RunInputDTO{Url: "https://github.com/LuisGaravaso/goexpert-auction", Requests: 5, Concurrency: 1, ShowData: false}
	ctx := context.Background()

	// Act
	output, err := uc.Run(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, input.Url, output.Url)
	assert.Equal(t, input.Requests, output.Requests)
	assert.Equal(t, input.Concurrency, output.Concurrency)
	assert.NotEmpty(t, output.Id)
	assert.NotEmpty(t, output.TimestampStart)
	assert.NotEmpty(t, output.TimestampEnd)
	assert.GreaterOrEqual(t, output.TestDurationInSeconds, 0)
	assert.Len(t, output.Data, 0)
	t.Log(output)
}

func Test_MustMakeRequestsConcurrently(t *testing.T) {
	// Arrange
	repo := &repository.MockRepository{}
	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	uc := run.NewRunUseCase(repo)
	requests := 9
	concurrency := 3
	showData := false
	input := run.RunInputDTO{Url: "https://github.com/LuisGaravaso/goexpert-auction", Requests: requests, Concurrency: concurrency, ShowData: showData}
	ctx := context.Background()

	// Act
	output, err := uc.Run(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, input.Url, output.Url)
	assert.Equal(t, input.Requests, output.Requests)
	assert.Equal(t, input.Concurrency, output.Concurrency)
	assert.NotEmpty(t, output.Id)
	assert.NotEmpty(t, output.TimestampStart)
	assert.NotEmpty(t, output.TimestampEnd)
	assert.GreaterOrEqual(t, output.TestDurationInSeconds, 0)
	if showData {
		assert.Len(t, output.Data, requests)
	}
	t.Log(output)
}
