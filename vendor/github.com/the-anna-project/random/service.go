// Package random provides a service implementation creating random numbers.
package random

import (
	"crypto/rand"
	"io"
	"math/big"
	"time"

	"github.com/cenk/backoff"
)

// ServiceConfig represents the configuration used to create a new random
// service.
type ServiceConfig struct {
	// Dependencies.

	// BackoffFactory is supposed to be able to create a new spec.Backoff. Retry
	// implementations can make use of this to decide when to retry.
	BackoffFactory func() Backoff
	// RandFactory represents a service returning random values. Here e.g.
	// crypto/rand.Int can be used.
	RandFactory func(rand io.Reader, max *big.Int) (n *big.Int, err error)

	// Settings.

	// RandReader represents an instance of a cryptographically strong
	// pseudo-random generator. Here e.g. crypto/rand.Reader can be used.
	RandReader io.Reader
	// Timeout represents the deadline being waited during random number creation
	// before returning a timeout error.
	Timeout time.Duration
}

// DefaultServiceConfig provides a default configuration to create a new random
// service by best effort.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		// Dependencies.
		BackoffFactory: func() Backoff {
			return &backoff.StopBackOff{}
		},
		RandFactory: rand.Int,

		// Settings.
		RandReader: rand.Reader,
		Timeout:    1 * time.Second,
	}
}

// NewService creates a new configured random service.
func NewService(config ServiceConfig) (Service, error) {
	// Dependencies.
	if config.BackoffFactory == nil {
		return nil, maskAnyf(invalidConfigError, "backoff factory must not be empty")
	}
	if config.RandFactory == nil {
		return nil, maskAnyf(invalidConfigError, "rand factory must not be empty")
	}

	// Settings.
	if config.RandReader == nil {
		return nil, maskAnyf(invalidConfigError, "rand reader must not be empty")
	}

	newService := &service{
		// Dependencies.
		backoffFactory: config.BackoffFactory,
		randFactory:    config.RandFactory,

		// Settings.
		randReader: config.RandReader,
		timeout:    config.Timeout,
	}

	return newService, nil
}

type service struct {
	// Dependencies.
	backoffFactory func() Backoff
	randFactory    func(rand io.Reader, max *big.Int) (n *big.Int, err error)

	// Settings.
	randReader io.Reader
	timeout    time.Duration
}

func (s *service) CreateMax(max int) (int, error) {
	// Define the action.
	var result int
	action := func() error {
		done := make(chan struct{}, 1)
		fail := make(chan error, 1)

		go func() {
			m := big.NewInt(int64(max))
			j, err := s.randFactory(s.randReader, m)
			if err != nil {
				fail <- maskAny(err)
				return
			}

			result = int(j.Int64())

			done <- struct{}{}
		}()

		select {
		case <-time.After(s.timeout):
			return maskAnyf(timeoutError, "after %s", s.timeout)
		case err := <-fail:
			return maskAny(err)
		case <-done:
			return nil
		}
	}

	// Execute the action wrapped with a retrier.
	err := backoff.Retry(action, s.backoffFactory())
	if err != nil {
		return 0, maskAny(err)
	}

	return result, nil
}

func (s *service) CreateNMax(n, max int) ([]int, error) {
	var result []int

	for i := 0; i < n; i++ {
		j, err := s.CreateMax(max)
		if err != nil {
			return nil, maskAny(err)
		}

		result = append(result, j)
	}

	return result, nil
}
