package app

import (
	"context"
	"log/slog"
	"os"
	"syscall"
	"time"
)

type options struct {
	ctx          context.Context
	sigs         []os.Signal
	logger       *slog.Logger
	startTimeout time.Duration
	stopTimeout  time.Duration
}

func defaultOptions() options {
	return options{
		ctx:          context.Background(),
		sigs:         []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT},
		logger:       slog.Default(),
		startTimeout: 30 * time.Second,
		stopTimeout:  30 * time.Second,
	}
}

// Option is a function that configures the App.
type Option func(*options)

func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func WithSignals(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

func WithLogger(logger *slog.Logger) Option {
	return func(o *options) { o.logger = logger }
}

func WithStartTimeout(d time.Duration) Option {
	return func(o *options) { o.startTimeout = d }
}

func WithStopTimeout(d time.Duration) Option {
	return func(o *options) { o.stopTimeout = d }
}
