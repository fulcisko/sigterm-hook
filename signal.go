package sigtermhook

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// DefaultShutdownSignals are the OS signals that trigger graceful shutdown.
var DefaultShutdownSignals = []os.Signal{
	syscall.SIGTERM,
	syscall.SIGINT,
}

// SignalWatcher watches for OS signals and triggers shutdown on the Manager.
type SignalWatcher struct {
	manager *Manager
	signals []os.Signal
}

// NewSignalWatcher creates a SignalWatcher that will call manager.Shutdown
// when any of the provided signals is received. If no signals are provided,
// DefaultShutdownSignals is used.
func NewSignalWatcher(manager *Manager, signals ...os.Signal) *SignalWatcher {
	if len(signals) == 0 {
		signals = DefaultShutdownSignals
	}
	return &SignalWatcher{
		manager: manager,
		signals: signals,
	}
}

// Watch blocks until one of the registered signals is received or the provided
// context is cancelled. It then initiates shutdown and returns the error
// returned by Manager.Shutdown.
func (sw *SignalWatcher) Watch(ctx context.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sw.signals...)
	defer signal.Stop(ch)

	select {
	case sig := <-ch:
		if sw.manager.logger != nil {
			sw.manager.logger.Printf("sigterm-hook: received signal %s, initiating shutdown", sig)
		}
		return sw.manager.Shutdown(ctx)
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WatchBackground starts Watch in a goroutine and returns a channel that
// receives the shutdown error (or nil) when the watcher exits.
func (sw *SignalWatcher) WatchBackground(ctx context.Context) <-chan error {
	result := make(chan error, 1)
	go func() {
		result <- sw.Watch(ctx)
	}()
	return result
}
