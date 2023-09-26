package main

import (
	"context"

	"github.com/spf13/cobra"

	"darvaza.org/sidecar/pkg/certmagic"
	"darvaza.org/sidecar/pkg/sidecar"
	"darvaza.org/slog"

	"darvaza.org/kube/pkg/edge"
)

// Application is what Service runs
type Application struct {
	srv   *sidecar.Server
	proxy *edge.Proxy
}

// Init initializes the [Application]
func (app *Application) Init(ctx context.Context, cfg *Config, logger slog.Logger) error {
	// TLS Store
	cm, err := cfg.TLS.New(
		certmagic.WithContext(ctx),
		certmagic.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	// Listener
	cfg.Server.Context = ctx
	cfg.Server.Logger = logger
	cfg.Server.Store = cm

	srv, err := cfg.Server.New()
	if err != nil {
		return err
	}

	// Handler
	cfg.Proxy.Context = ctx
	cfg.Proxy.Logger = logger
	cfg.Proxy.Store = cm

	r, err := cfg.Proxy.New()
	if err != nil {
		srv.Cancel()
		return err
	}

	app.srv = srv
	app.proxy = r
	return nil
}

// Run runs the application
func (app *Application) Run(_ context.Context, _ *cobra.Command, _ []string) error {
	// TODO: shutdown when the given context is cancelled
	return app.srv.ListenAndServe(app.proxy)
}
