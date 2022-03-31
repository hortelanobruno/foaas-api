package ratelimiter

import (
	"sync"
	"time"
)

type LocalRateLimiter struct {
	rateLimitCount           int
	rateWindowInMilliseconds time.Duration
	requestsByUser           map[string][]time.Time
	mutex                    *sync.Mutex
	now                      func() time.Time
}

func NewLocalRateLimiter(rateLimitCount int, rateWindowInMilliseconds time.Duration) *LocalRateLimiter {
	return &LocalRateLimiter{
		rateLimitCount:           rateLimitCount,
		rateWindowInMilliseconds: rateWindowInMilliseconds,
		requestsByUser:           make(map[string][]time.Time, 0),
		mutex:                    &sync.Mutex{},
		now:                      time.Now,
	}
}

// AllowRequest returns true if in the last past X milliseconds, there were fewer requests than the rate limit.
// Remove all the old requests from the map.
func (s *LocalRateLimiter) AllowRequest(userID string) bool {
	now := s.now()
	s.mutex.Lock()
	defer s.mutex.Unlock()

	requests, exists := s.requestsByUser[userID]
	if !exists {
		s.requestsByUser[userID] = []time.Time{now}
		return true
	}

	newRequests := s.getRequestsInTheWindowTime(requests, now)
	s.requestsByUser[userID] = newRequests
	if len(newRequests) >= s.rateLimitCount {
		return false
	}

	s.requestsByUser[userID] = append(newRequests, now)
	return true
}

func (s *LocalRateLimiter) getRequestsInTheWindowTime(requests []time.Time, now time.Time) []time.Time {
	newRequests := make([]time.Time, 0)
	for _, request := range requests {
		if now.Sub(request) <= s.rateWindowInMilliseconds {
			newRequests = append(newRequests, request)
		}
	}
	return newRequests
}
