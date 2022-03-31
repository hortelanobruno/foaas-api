package server

type Options struct {
	LogLevel                      string
	RateLimitEnable               bool
	RateLimitCount                int
	RateLimitWindowInMilliseconds int
	TimeoutInMilliseconds         int
}
