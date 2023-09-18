// Package main implements a Kubernetes Cloud Controller Manager
package main

import (
	"os"

	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"          // register optional JSON log format
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugins
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/klog/v2"

	"darvaza.org/kube/pkg/ccm"
)

func main() {
	cmd, err := ccm.NewCommand()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	code := cli.Run(cmd)
	os.Exit(code)
}
