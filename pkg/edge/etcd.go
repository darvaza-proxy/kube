package edge

import (
	client "go.etcd.io/etcd/client/v3"
)

// Export generates a config for the etcd client
func (*Config) Export() (client.Config, error) {
	return client.Config{}, nil
}

func newEtcdClient(cfg *Config) (*client.Client, error) {
	conf, err := cfg.Export()
	if err != nil {
		return nil, err
	}

	ec, err := client.New(conf)
	if err != nil {
		return nil, err
	}

	return ec, nil
}
