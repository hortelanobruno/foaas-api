package ratelimiter

type RateLimiter interface {
	AllowRequest(userId string) bool
}
