package edge

import (
	"net/http"

	etcdclient "go.etcd.io/etcd/client/v3"
)

var (
	_ http.Handler = (*Proxy)(nil)
)

// Proxy implements a reverse proxy for Kubernetes nodes
type Proxy struct {
	cfg Config

	ec *etcdclient.Client
}

func (p *Proxy) init() error {
	ec, err := newEtcdClient(&p.cfg)
	if err != nil {
		return err
	}

	p.ec = ec
	return nil
}

func (*Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	http.NotFound(rw, req)
}

// NewProxy creates a [Proxy] based on the given [Config]
func NewProxy(cfg Config) (*Proxy, error) {
	if err := cfg.Prepare(); err != nil {
		return nil, err
	}

	p := &Proxy{
		cfg: cfg,
	}

	if err := p.init(); err != nil {
		return nil, err
	}

	return p, nil
}
