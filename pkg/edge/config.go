package edge

import (
	"context"

	"darvaza.org/slog"
)

// Config describes the operation of the proxy
type Config struct {
	Logger  slog.Logger     `json:"-" yaml:"-" toml:"-"`
	Context context.Context `json:"-" yaml:"-" toml:"-"`
}

// Prepare fills any gap in the Config and validates its content
func (*Config) Prepare() error {
	return nil
}

// New creates a [Proxy] from the [Config]
func (cfg Config) New() (*Proxy, error) {
	return NewProxy(cfg)
}
