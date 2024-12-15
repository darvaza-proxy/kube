package ccm

import (
	"io"
	"os"

	"darvaza.org/core"

	cloudprovider "k8s.io/cloud-provider"
)

const (
	// ProviderName is the name used when registering ourselves
	// as cloud providers
	ProviderName = "darvaza"
)

var (
	_ cloudprovider.Interface = (*Provider)(nil)
)

// Provider is an instance of the cloud we implement
type Provider struct {
	configContent []byte
}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping or run custom controllers specific to the cloud provider.
// Any tasks started here should be cleaned up when the stop channel closes.
func (*Provider) Initialize(_ cloudprovider.ControllerClientBuilder, _ <-chan struct{}) {}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (*Provider) LoadBalancer() (cloudprovider.LoadBalancer, bool) { return nil, false }

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (*Provider) Instances() (cloudprovider.Instances, bool) { return nil, false }

// InstancesV2 is an implementation for instances and should only be implemented by external cloud providers.
// Implementing InstancesV2 is behaviorally identical to Instances but is optimized to significantly reduce
// API calls to the cloud provider when registering and syncing nodes. Implementation of this interface will
// disable calls to the Zones interface. Also returns true if the interface is supported, false otherwise.
func (*Provider) InstancesV2() (cloudprovider.InstancesV2, bool) { return nil, false }

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
// DEPRECATED: Zones is deprecated in favor of retrieving zone/region information from InstancesV2.
// This interface will not be called if InstancesV2 is enabled.
func (*Provider) Zones() (cloudprovider.Zones, bool) { return nil, false }

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (*Provider) Clusters() (cloudprovider.Clusters, bool) { return nil, false }

// Routes returns a routes interface along with whether the interface is supported.
func (*Provider) Routes() (cloudprovider.Routes, bool) { return nil, false }

// ProviderName returns the cloud provider ID.
func (*Provider) ProviderName() string { return ProviderName }

// HasClusterID returns true if a ClusterID is required and set
func (*Provider) HasClusterID() bool { return false }

func providerFactory(r io.Reader) (cloudprovider.Interface, error) {
	b, err := io.ReadAll(r)
	switch {
	case os.IsNotExist(err):
		// no config file, fine for now
	case err != nil:
		// another error, fatal
		return nil, core.Wrapf(err, "cloud-provider-%s: failed to read config", ProviderName)
	}

	p := &Provider{
		configContent: b,
	}

	return p, nil
}

func init() {
	// register provider
	cloudprovider.RegisterCloudProvider(ProviderName, providerFactory)
}
