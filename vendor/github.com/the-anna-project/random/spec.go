package random

import "time"

// Backoff represents the object managing backoff algorithms to retry actions.
type Backoff interface {
	// NextBackOff provides the duration expected to wait before retrying an
	// action. time.Duration = -1 indicates that no more retry should be
	// attempted.
	NextBackOff() time.Duration
	// Reset sets the backoff back to its initial state.
	Reset()
}

// Service creates pseudo random numbers. The service might implement retries
// using backoff strategies and timeouts.
type Service interface {
	// CreateMax tries to create one new pseudo random number. The generated
	// number is within the range [0 max), which means that max is exclusive.
	CreateMax(max int) (int, error)
	// CreateNMax tries to create a list of new pseudo random numbers. n
	// represents the number of pseudo random numbers in the returned list. The
	// generated numbers are within the range [0 max), which means that max is
	// exclusive.
	CreateNMax(n, max int) ([]int, error)
}
