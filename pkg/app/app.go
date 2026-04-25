package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

// Hook is a pair of start and stop callbacks.
type Hook struct {
	Name    string
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

// App manages the lifecycle of a Go application.
type App struct {
	opts options
	ctx    context.Context
	cancel func()
	hooks  []Hook
	mu     sync.Mutex
}

// New creates a new application manager with the given options.
func New(opts ...Option) *App {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{opts: o, ctx: ctx, cancel: cancel}
}

// AppendHook adds a lifecycle hook to the application.
func (a *App) AppendHook(h Hook) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.hooks = append(a.hooks, h)
}

// Run starts the application and waits for an interrupt signal.
func (a *App) Run() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, a.opts.sigs...)

	a.opts.logger.Info("Starting application...")
	startCtx, cancelStart := context.WithTimeout(a.ctx, a.opts.startTimeout)
	defer cancelStart()

	if err := a.start(startCtx); err != nil {
		a.opts.logger.Error("Failed to start application", "error", err)
		_ = a.Stop()
		return err
	}
	a.opts.logger.Info("Application started successfully")

	select {
	case sig := <-stop:
		a.opts.logger.Info("Received stop signal", "signal", sig.String())
	case <-a.ctx.Done():
		a.opts.logger.Info("Application context done")
	}

	return a.Stop()
}

// Stop gracefully shuts down the application.
func (a *App) Stop() error {
	a.opts.logger.Info("Stopping application...")
	a.cancel()

	stopCtx, cancelStop := context.WithTimeout(context.Background(), a.opts.stopTimeout)
	defer cancelStop()

	var errs []error
	a.mu.Lock()
	hooks := make([]Hook, len(a.hooks))
	copy(hooks, a.hooks)
	a.mu.Unlock()

	for i := len(hooks) - 1; i >= 0; i-- {
		hook := hooks[i]
		if hook.OnStop != nil {
			hookName := hook.Name
			if hookName == "" {
				hookName = fmt.Sprintf("hook-%d", i)
			}
			a.opts.logger.Debug("Stopping component", "name", hookName)
			if err := hook.OnStop(stopCtx); err != nil {
				a.opts.logger.Error("Failed to stop component", "name", hookName, "error", err)
				errs = append(errs, err)
			} else {
				a.opts.logger.Debug("Successfully stopped component", "name", hookName)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("application stopped with errors")
	}
	a.opts.logger.Info("Application stopped successfully")
	return nil
}

func (a *App) start(ctx context.Context) error {
	a.mu.Lock()
	hooks := make([]Hook, len(a.hooks))
	copy(hooks, a.hooks)
	a.mu.Unlock()

	for i, hook := range hooks {
		if hook.OnStart != nil {
			hookName := hook.Name
			if hookName == "" {
				hookName = fmt.Sprintf("hook-%d", i)
			}
			a.opts.logger.Debug("Starting component", "name", hookName)
			if err := hook.OnStart(ctx); err != nil {
				return fmt.Errorf("failed to start component %s: %w", hookName, err)
			}
			a.opts.logger.Debug("Successfully started component", "name", hookName)
		}
	}
	return nil
}
