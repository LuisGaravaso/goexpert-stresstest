package entity

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

const (
	ErrInvalidURL             = "invalid url, must be in the format http://example.com or https://example.com"
	ErrNonNegativeRequests    = "requests must be greater than zero"
	ErrNonNegativeConcurrency = "concurrency must be greater than zero"
)

type TestRun struct {
	Id          string
	Url         string
	Requests    int
	Concurrency int
	Timestamp   time.Time
}

type TestRunOptions struct {
	Requests    int
	Concurrency int
}

func NewTestRun(url string, opts *TestRunOptions) (*TestRun, error) {
	requests := 100   // default
	concurrency := 10 // default
	if opts != nil {
		if opts.Requests != 0 {
			requests = opts.Requests
		}
		if opts.Concurrency != 0 {
			concurrency = opts.Concurrency
		}
	}
	if concurrency > requests {
		concurrency = requests
	}

	tr := &TestRun{
		Id:          uuid.New().String(),
		Url:         url,
		Requests:    requests,
		Concurrency: concurrency,
		Timestamp:   time.Now(),
	}

	if err := tr.Validate(); err != nil {
		return nil, err
	}

	return tr, nil

}

func (tr *TestRun) Validate() error {
	if !IsValidURL(tr.Url) {
		return errors.New(ErrInvalidURL)
	}
	if tr.Requests <= 0 {
		return errors.New(ErrNonNegativeRequests)
	}
	if tr.Concurrency <= 0 {
		return errors.New(ErrNonNegativeConcurrency)
	}
	return nil
}

func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
