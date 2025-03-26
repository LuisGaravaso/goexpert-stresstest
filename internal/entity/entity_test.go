package entity_test

import (
	"stresstest/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestRun_ValidParams(t *testing.T) {
	opts := &entity.TestRunOptions{Requests: 50, Concurrency: 10}
	tr, err := entity.NewTestRun("http://example.com", opts)

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", tr.Url)
	assert.Equal(t, 50, tr.Requests)
	assert.Equal(t, 10, tr.Concurrency)
	assert.NotEmpty(t, tr.Id)
}

func TestNewTestRun_InvalidURL(t *testing.T) {
	opts := &entity.TestRunOptions{Requests: 10, Concurrency: 5}
	tr, err := entity.NewTestRun("invalid-url", opts)

	assert.Nil(t, tr)
	assert.EqualError(t, err, entity.ErrInvalidURL)
}

func TestNewTestRun_NegativeRequests(t *testing.T) {
	opts := &entity.TestRunOptions{Requests: -1, Concurrency: 5}
	tr, err := entity.NewTestRun("http://example.com", opts)

	assert.Nil(t, tr)
	assert.EqualError(t, err, entity.ErrNonNegativeRequests)
}

func TestNewTestRun_NegativeConcurrency(t *testing.T) {
	opts := &entity.TestRunOptions{Requests: 10, Concurrency: -1}
	tr, err := entity.NewTestRun("http://example.com", opts)

	assert.Nil(t, tr)
	assert.EqualError(t, err, entity.ErrNonNegativeConcurrency)
}

func TestNewTestRun_ConcurrencyGreaterThanRequests(t *testing.T) {
	opts := &entity.TestRunOptions{Requests: 10, Concurrency: 20}
	tr, err := entity.NewTestRun("http://example.com", opts)

	assert.NoError(t, err)
	assert.Equal(t, 10, tr.Requests)
	assert.Equal(t, 10, tr.Concurrency) // ajustado automaticamente
}

func TestNewTestRun_DefaultValues(t *testing.T) {
	tr, err := entity.NewTestRun("http://example.com", nil)

	assert.NoError(t, err)
	assert.Equal(t, 100, tr.Requests)
	assert.Equal(t, 10, tr.Concurrency)
}
