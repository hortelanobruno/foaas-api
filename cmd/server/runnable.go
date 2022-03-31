package server

import (
	"github.com/hortelanobruno/foaas-api/domain/service"
	"github.com/hortelanobruno/foaas-api/domain/service/handler"
	"github.com/hortelanobruno/foaas-api/domain/validator"
	"github.com/hortelanobruno/foaas-api/http"
	"github.com/hortelanobruno/foaas-api/ratelimiter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

type Runnable struct{}

func NewRunnable() *Runnable {
	return &Runnable{}
}

func (r *Runnable) Cmd() *cobra.Command {
	options := &Options{}

	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Runs foaas API",
		Long:  `Runs foaas API`,
	}

	cmd.Flags().StringVar(&options.LogLevel, "log-level", defaultLogLevel, "log leve to use")
	cmd.Flags().BoolVar(&options.RateLimitEnable, "rate-limit-enable", defaultRateLimitEnable, "switch to enable rate limiter")
	cmd.Flags().IntVar(&options.RateLimitCount, "rate-limit-count", defaultRateLimitCount, "maximum quantity of requests "+
		"that a user can do in a window of time")
	cmd.Flags().IntVar(&options.RateLimitWindowInMilliseconds, "rate-limit-window-in-milliseconds", defaultRateLimitWindowInMilliseconds,
		"window of time in milliseconds to limit the quantity of requests that a user can do")
	cmd.Flags().IntVar(&options.TimeoutInMilliseconds, "timeout-in-milliseconds", defaultTimeoutInMilliseconds,
		"timeout of the api calls")

	cmd.Run = func(_ *cobra.Command, _ []string) {
		server := r.Run(options)
		server.Start(defaultPort)
	}
	return cmd
}

func (r *Runnable) Run(options *Options) *Server {
	r.configureLog(options.LogLevel)

	var rateLimiter ratelimiter.RateLimiter
	if options.RateLimitEnable {
		rateLimiter = ratelimiter.NewLocalRateLimiter(
			options.RateLimitCount,
			time.Duration(options.RateLimitWindowInMilliseconds)*time.Millisecond)
	}

	httpClient := http.NewClientImpl(time.Duration(options.TimeoutInMilliseconds) * time.Millisecond)

	messageService := service.NewMessageServiceImpl(httpClient)
	messageValidator := validator.NewMessageValidatorImpl()
	messageHandler := handler.NewMessageHandler(messageValidator, messageService)

	return NewServer(messageHandler, rateLimiter)
}

func (r *Runnable) configureLog(logLevel string) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Warnf("Error passing the log level: %s", logLevel)
		return
	}

	logrus.Infof("Setting log level: %s", lvl.String())
	logrus.SetLevel(lvl)
}
