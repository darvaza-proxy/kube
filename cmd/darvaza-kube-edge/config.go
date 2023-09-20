package main

import (
	"os"

	"darvaza.org/sidecar/pkg/config"
	"darvaza.org/sidecar/pkg/sidecar"
	"github.com/spf13/cobra"

	"darvaza.org/kube/pkg/edge"
)

// Config is the configuration structure of this
// reverse proxy server
type Config struct {
	Server sidecar.Config `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Proxy  edge.Config    `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
}

// NewConfig creates a new [Config] initialized with
// default values
func NewConfig() *Config {
	cfg := new(Config)
	if err := config.Prepare(cfg); err != nil {
		// WTF
		fatal(err, "failed to initialize default config")
	}

	return cfg
}

// NewConfigFromFile loads the server's configuration from a file
// by name, expanding environment variables, filling gaps and
// validating its content.
func NewConfigFromFile(filename string) (*Config, error) {
	cfg := new(Config)
	err := config.LoadFile(filename, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// ReadInFile loads the server's configuration from a TOML file
// by name, expanding environment variables, filling gaps and
// validating its content. On error the object isn't touched.
func (cfg *Config) ReadInFile(filename string) error {
	c, err := NewConfigFromFile(filename)
	if err != nil {
		return err
	}

	if cfg != nil {
		*cfg = *c
	}
	return nil
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump provides a text representation of the loaded configuration",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, err := newConfigTryFile(cfgFile)
		if err != nil {
			return err
		}

		format, err := cmd.Flags().GetString("output-format")
		if err != nil {
			return err
		}

		enc, err := config.NewEncoder(format)
		if err != nil {
			return err
		}

		_, err = enc.WriteTo(cfg, os.Stdout)
		return err
	},
}

func init() {
	dumpCmd.Flags().String("output-format", "toml", "toml, yaml or json")

	svc.AddCommand(dumpCmd)
}
