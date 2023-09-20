// Package main implements the darvaza edge server for kubernetes nodes
package main

import (
	"context"
	"os"

	"darvaza.org/sidecar/pkg/service"
	"darvaza.org/slog"
	"github.com/spf13/cobra"

	"darvaza.org/kube/pkg/version"
)

const (
	// CmdName is the name of the executable
	CmdName = "darvaza-kube-edge"

	// DefaultConfigFile is the name of the configuration
	// file to be read if none is specified
	DefaultConfigFile = CmdName + ".toml"
)

var (
	app     Application
	cfgFile string
)

var (
	svcConfig = &service.Config{
		Name:        CmdName,
		Description: "Darvaza Edge Server for Kubernetes nodes",
		Short:       "runs the proxy",
		Version:     version.Version,

		Prepare: func(ctx context.Context, _ *cobra.Command, _ []string) error {
			cfg, err := newConfigTryFile(cfgFile)
			if err != nil {
				return err
			}

			return app.Init(ctx, cfg, newLogger())
		},
		Run: app.Run,
	}

	svc = service.Must(svcConfig)
)

func newConfigTryFile(filename string) (*Config, error) {
	if filename != "" {
		cfg, err := NewConfigFromFile(filename)
		switch {
		case err == nil:
			// loaded
			return cfg, nil
		case os.IsNotExist(err) && filename == DefaultConfigFile:
			// missing DefaultConfigFile, ignore
		default:
			// bad config file
			return nil, err
		}
	}

	// didn't load, create one
	return NewConfig(), nil
}

// main invokes cobra
func main() {
	err := svc.Execute()
	code, err := service.AsExitStatus(err)

	if err != nil {
		newLogger().Error().
			WithField(slog.ErrorFieldName, err).
			Print()
	}

	os.Exit(code)
}

func init() {
	pflags := svc.PersistentFlags()
	pflags.StringVarP(&cfgFile, "config-file", "f", DefaultConfigFile, "config file (TOML format)")
}
