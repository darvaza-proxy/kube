// Package main implements a Kubernetes Cloud Controller Manager
package main

import (
	"os"

	"k8s.io/component-base/cli"
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
