package ccm

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"

	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	"k8s.io/cloud-provider/app/config"
	"k8s.io/cloud-provider/names"
	"k8s.io/cloud-provider/options"
)

// Options describes the Cloud Controller Manager (CCM) we want to
// build
type Options struct {
	ccmOptions   *options.CloudControllerManagerOptions
	initializers map[string]app.ControllerInitFuncConstructor
	aliases      map[string]string
	fss          flag.NamedFlagSets
	cloud        cloudprovider.Interface
}

// FlagSets returns the cli flags
func (opts *Options) FlagSets() *flag.NamedFlagSets {
	return &opts.fss
}

// FlagSet returns a cli flagset by name
func (opts *Options) FlagSet(name string) *pflag.FlagSet {
	return opts.fss.FlagSet(name)
}

// Cloud returns the initialized cloud provider
func (opts *Options) Cloud() cloudprovider.Interface {
	return opts.cloud
}

func (opts *Options) cloudInitializer(c *config.CompletedConfig) cloudprovider.Interface {
	cc := c.ComponentConfig.KubeCloudShared.CloudProvider

	cloud, err := cloudprovider.InitCloudProvider(cc.Name, cc.CloudConfigFile)
	switch {
	case err != nil:
		klog.Fatalf("cloud-provider-%s: failed to initialize: %v", cc.Name, err)
	case cloud == nil:
		klog.Fatalf("cloud-provider-%s: nil value", cc.Name)
	case cloud.HasClusterID():
		// ready
	case c.ComponentConfig.KubeCloudShared.AllowUntaggedCloud:
		// untagged but allowed
		klog.Warning("cloud-provider-%s: TODO: a ClusterID will be required in the future", cc.Name)
	default:
		// untagged not allowed
		klog.Fatalf("cloud-provider-%s: no ClusterID found", cc.Name)
	}

	// remember
	opts.cloud = cloud

	return cloud
}

// NewCommand creates a [cobra.Command] based on the [Options]
func (opts *Options) NewCommand() (*cobra.Command, error) {
	cmd := app.NewCloudControllerManagerCommand(opts.ccmOptions,
		opts.cloudInitializer,
		opts.initializers,
		opts.aliases,
		opts.fss,
		wait.NeverStop)
	return cmd, nil
}

// NewOptions create Options to build a Cloud Controller Manager instance
func NewOptions() (*Options, error) {
	ccmOptions, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		return nil, err
	}

	opts := &Options{
		ccmOptions:   ccmOptions,
		initializers: app.DefaultInitFuncConstructors,
		aliases:      names.CCMControllerAliases(),
	}

	return opts, nil
}

// NewCommand creates a [cobra.Command] for the Cloud Controller Manager
// using the default options
func NewCommand() (*cobra.Command, error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return opts.NewCommand()
}
