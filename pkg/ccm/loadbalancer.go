package ccm

// revive:disable:line-length-limit
// revive:disable:redundant-import-alias

import (
	"context"
	"fmt"
	"sync"

	v1 "k8s.io/api/core/v1"

	cloudprovider "k8s.io/cloud-provider"
)

var (
	_ cloudprovider.LoadBalancer = (*LoadBalancerManager)(nil)
)

// LoadBalancer is a particular load balancer for a specific service
type LoadBalancer struct{}

// Status returns the status of the load balancer
func (LoadBalancer) Status() (status *v1.LoadBalancerStatus, err error) {
	return nil, ErrNotImplemented
}

func (LoadBalancer) getStatus() *v1.LoadBalancerStatus {
	return nil
}

func (LoadBalancer) updateNodes(_ context.Context, _ []*v1.Node) (err error) {
	return ErrNotImplemented
}

// LoadBalancerManager implements a LoadBalancer Manager
type LoadBalancerManager struct {
	mu sync.Mutex
	p  *Provider
}

func (lbm *LoadBalancerManager) init(p *Provider) error {
	if p == nil {
		panic("unreachable")
	}

	lbm.p = p
	return nil
}

func (*LoadBalancerManager) getName(clusterName string, service *v1.Service) string {
	return fmt.Sprintf("%s:%s:%s", "darvaza", clusterName, service.Name)
}

func (*LoadBalancerManager) getBalancer(_ context.Context, _ string) (*LoadBalancer, bool, error) {
	// Implementations may return a (possibly wrapped) api.RetryError to enforce
	// backing off at a fixed duration. This can be used for cases like when the
	// load balancer is not ready yet (e.g., it is still being provisioned) and
	// polling at a fixed rate is preferred over backing off exponentially in
	// order to minimize latency.
	return nil, false, ErrNotImplemented
}

func (*LoadBalancerManager) newBalancer(_ context.Context, _, _ string, _ *v1.Service) (*LoadBalancer, error) {
	return nil, ErrNotImplemented
}

func (*LoadBalancerManager) deleteBalancer(_ context.Context, _ string, _ *LoadBalancer) error {
	return ErrNotImplemented
}

// GetLoadBalancer returns whether the specified load balancer exists, and
// if so, what its status is.
func (lbm *LoadBalancerManager) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	lbm.mu.Lock()
	defer lbm.mu.Unlock()

	name := lbm.getName(clusterName, service)
	lb, exists, err := lbm.getBalancer(ctx, name)
	if !exists || err == nil {
		return nil, exists, err
	}

	status, err = lb.Status()
	switch {
	case err != nil:
		return nil, true, err
	default:
		return status, true, nil
	}
}

// GetLoadBalancerName returns the name of the load balancer.
func (lbm *LoadBalancerManager) GetLoadBalancerName(_ context.Context, clusterName string, service *v1.Service) string {
	lbm.mu.Lock()
	defer lbm.mu.Unlock()

	return lbm.getName(clusterName, service)
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
func (lbm *LoadBalancerManager) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	lbm.mu.Lock()
	defer lbm.mu.Unlock()

	name := lbm.getName(clusterName, service)
	lb, exists, err := lbm.getBalancer(ctx, name)

	if !exists && err == nil {
		lb, err = lbm.newBalancer(ctx, name, clusterName, service)
	}

	if err != nil {
		// failed to get or create
		return nil, err
	} else if err = lb.updateNodes(ctx, nodes); err != nil {
		// failed to update nodes
		return nil, err
	}

	return lb.getStatus(), nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
func (lbm *LoadBalancerManager) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	lbm.mu.Lock()
	defer lbm.mu.Unlock()

	name := lbm.getName(clusterName, service)
	lb, exists, err := lbm.getBalancer(ctx, name)
	if exists && err == nil {
		// update
		err = lb.updateNodes(ctx, nodes)
	}

	return err
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
func (lbm *LoadBalancerManager) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	lbm.mu.Lock()
	defer lbm.mu.Unlock()

	name := lbm.getName(clusterName, service)
	lb, exists, err := lbm.getBalancer(ctx, name)
	if !exists || err != nil {
		return err
	}

	return lbm.deleteBalancer(ctx, name, lb)
}
