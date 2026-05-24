package sigterm

import "time"

// TimeoutConfig holds timeout settings for the shutdown manager.
type TimeoutConfig struct {
	// GlobalTimeout is the maximum time allowed for the entire shutdown sequence.
	// Zero means no global timeout.
	GlobalTimeout time.Duration

	// HandlerTimeout is the default maximum time allowed per handler.
	// Zero means no per-handler timeout.
	HandlerTimeout time.Duration

	// DrainTimeout is the maximum time to wait for in-flight requests to drain
	// before proceeding with handler teardown.
	// Zero means no drain timeout.
	DrainTimeout time.Duration
}

// Validate checks that timeout values are non-negative.
func (c TimeoutConfig) Validate() error {
	if c.GlobalTimeout < 0 {
		return ErrInvalidTimeout{Field: "GlobalTimeout", Value: c.GlobalTimeout}
	}
	if c.HandlerTimeout < 0 {
		return ErrInvalidTimeout{Field: "HandlerTimeout", Value: c.HandlerTimeout}
	}
	if c.DrainTimeout < 0 {
		return ErrInvalidTimeout{Field: "DrainTimeout", Value: c.DrainTimeout}
	}
	return nil
}

// ErrInvalidTimeout is returned when a timeout value is invalid.
type ErrInvalidTimeout struct {
	Field string
	Value time.Duration
}

func (e ErrInvalidTimeout) Error() string {
	return "sigterm: invalid timeout for " + e.Field + ": " + e.Value.String()
}

// mergeTimeoutConfig merges a per-handler timeout with the global config.
// The per-handler timeout takes precedence if set; otherwise the global
// HandlerTimeout is used.
func mergeTimeoutConfig(global TimeoutConfig, perHandler time.Duration) time.Duration {
	if perHandler > 0 {
		return perHandler
	}
	return global.HandlerTimeout
}
