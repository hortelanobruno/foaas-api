package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetRequestsInTheWindowTime(t *testing.T) {
	cases := []struct {
		name                          string
		inputRequests                 []time.Time
		inputNow                      time.Time
		inputRateWindowInMilliseconds int
		expectedOutput                []time.Time
	}{
		{
			"Should return empty when requests are nil",
			nil,
			time.Now(),
			10,
			[]time.Time{},
		},
		{
			"Should return empty when requests are empty",
			[]time.Time{},
			time.Now(),
			10,
			[]time.Time{},
		},
		{
			"Should return empty when all the requests were on the past",
			[]time.Time{
				time.Date(2022, time.March, 30, 0, 0, 0, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 0, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 1, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 2, 00, time.UTC),
			},
			time.Date(2022, time.March, 30, 0, 0, 5, 00, time.UTC),
			1000,
			[]time.Time{},
		},
		{
			"Should return the last 2 requests when they are inside of the window time",
			[]time.Time{
				time.Date(2022, time.March, 30, 0, 0, 0, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 1, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 2, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 4, 00, time.UTC),
			},
			time.Date(2022, time.March, 30, 0, 0, 5, 00, time.UTC),
			1000 * 4,
			[]time.Time{
				time.Date(2022, time.March, 30, 0, 0, 1, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 2, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 4, 00, time.UTC),
			},
		},
		{
			"Should return the all the requests because they are inside the window time",
			[]time.Time{
				time.Date(2022, time.March, 30, 0, 0, 5, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 11, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 15, 00, time.UTC),
			},
			time.Date(2022, time.March, 30, 0, 0, 15, 00, time.UTC),
			1000 * 10,
			[]time.Time{
				time.Date(2022, time.March, 30, 0, 0, 5, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 11, 00, time.UTC),
				time.Date(2022, time.March, 30, 0, 0, 15, 00, time.UTC),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			rateLimiter := NewLocalRateLimiter(0,
				time.Duration(c.inputRateWindowInMilliseconds)*time.Millisecond)

			// Operation
			requests := rateLimiter.getRequestsInTheWindowTime(c.inputRequests, c.inputNow)

			// Validation
			assert.EqualValues(t, c.expectedOutput, requests)
		})
	}
}

func TestAllowRequestShouldReturnTrueWhenItIsTheFirstRequest(t *testing.T) {
	// Initialization
	userID := "123"

	rateLimiter := NewLocalRateLimiter(5, time.Duration(10000)*time.Millisecond)
	rateLimiter.now = func() time.Time {
		return time.Date(2022, time.March, 30, 0, 0, 0, 00, time.UTC)
	}

	// Operation
	isAllowed := rateLimiter.AllowRequest(userID)

	// Validation
	assert.True(t, isAllowed)
	assert.Len(t, rateLimiter.requestsByUser, 1)
	assert.EqualValues(t, rateLimiter.requestsByUser[userID],
		[]time.Time{rateLimiter.now()})
}

func TestAllowRequestShouldReturnTrueWhenRequestsInWindowTimeIsLessThanLimit(t *testing.T) {
	// Initialization
	userID := "123"

	rateLimiter := NewLocalRateLimiter(5, time.Duration(10000)*time.Millisecond)
	rateLimiter.now = func() time.Time {
		return time.Date(2022, time.March, 30, 0, 0, 18, 00, time.UTC)
	}

	rateLimiter.requestsByUser[userID] = []time.Time{
		time.Date(2022, time.March, 30, 0, 0, 1, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 2, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 13, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 14, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 17, 00, time.UTC),
	}

	// Operation
	isAllowed := rateLimiter.AllowRequest(userID)

	// Validation
	assert.True(t, isAllowed)
	assert.Len(t, rateLimiter.requestsByUser, 1)
	assert.EqualValues(t, rateLimiter.requestsByUser[userID],
		[]time.Time{
			time.Date(2022, time.March, 30, 0, 0, 13, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 14, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 17, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 18, 00, time.UTC),
		})
}

func TestAllowRequestShouldReturnFalseWhenRequestsInWindowTimeIsEqualThanLimit(t *testing.T) {
	// Initialization
	userID := "123"

	rateLimiter := NewLocalRateLimiter(3, time.Duration(10000)*time.Millisecond)
	rateLimiter.now = func() time.Time {
		return time.Date(2022, time.March, 30, 0, 0, 18, 00, time.UTC)
	}

	rateLimiter.requestsByUser[userID] = []time.Time{
		time.Date(2022, time.March, 30, 0, 0, 1, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 2, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 13, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 14, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 17, 00, time.UTC),
	}

	// Operation
	isAllowed := rateLimiter.AllowRequest(userID)

	// Validation
	assert.False(t, isAllowed)
	assert.Len(t, rateLimiter.requestsByUser, 1)
	assert.EqualValues(t, rateLimiter.requestsByUser[userID],
		[]time.Time{
			time.Date(2022, time.March, 30, 0, 0, 13, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 14, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 17, 00, time.UTC),
		})
}

func TestAllowRequestShouldReturnFalseWhenRequestsInWindowTimeIsGreaterThanLimit(t *testing.T) {
	// Initialization
	userID := "123"

	rateLimiter := NewLocalRateLimiter(3, time.Duration(10000)*time.Millisecond)
	rateLimiter.now = func() time.Time {
		return time.Date(2022, time.March, 30, 0, 0, 18, 00, time.UTC)
	}

	rateLimiter.requestsByUser[userID] = []time.Time{
		time.Date(2022, time.March, 30, 0, 0, 1, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 2, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 13, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 14, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 15, 00, time.UTC),
		time.Date(2022, time.March, 30, 0, 0, 17, 00, time.UTC),
	}

	// Operation
	isAllowed := rateLimiter.AllowRequest(userID)

	// Validation
	assert.False(t, isAllowed)
	assert.Len(t, rateLimiter.requestsByUser, 1)
	assert.EqualValues(t, rateLimiter.requestsByUser[userID],
		[]time.Time{
			time.Date(2022, time.March, 30, 0, 0, 13, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 14, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 15, 00, time.UTC),
			time.Date(2022, time.March, 30, 0, 0, 17, 00, time.UTC),
		})
}
