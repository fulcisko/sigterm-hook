package sigterm

import (
	"context"
	"fmt"
	"time"
)

// RetryConfig holds configuration for handler retry behavior.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// Delay is the wait duration between attempts.
	Delay time.Duration
	// OnRetry is an optional callback invoked before each retry.
	OnRetry func(attempt int, err error)
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       100 * time.Millisecond,
	}
}

// wrapWithRetry wraps a ShutdownFunc so that it is retried up to
// cfg.MaxAttempts times on failure. The context is respected between
// retries; if it is cancelled the function returns immediately.
func wrapWithRetry(fn ShutdownFunc, cfg RetryConfig) ShutdownFunc {
	if cfg.MaxAttempts <= 1 {
		return fn
	}
	return func(ctx context.Context) error {
		var lastErr error
		for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("context cancelled before attempt %d: %w", attempt, err)
			}
			lastErr = fn(ctx)
			if lastErr == nil {
				return nil
			}
			if attempt < cfg.MaxAttempts {
				if cfg.OnRetry != nil {
					cfg.OnRetry(attempt, lastErr)
				}
				select {
				case <-time.After(cfg.Delay):
				case <-ctx.Done():
					return fmt.Errorf("context cancelled during retry delay: %w", ctx.Err())
				}
			}
		}
		return fmt.Errorf("handler failed after %d attempts: %w", cfg.MaxAttempts, lastErr)
	}
}
