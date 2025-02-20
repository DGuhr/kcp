/*
Copyright 2021 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	goflags "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"

	crdcmd "github.com/kcp-dev/kcp/pkg/cliplugins/crd/cmd"
	workloadcmd "github.com/kcp-dev/kcp/pkg/cliplugins/workload/cmd"
	workspacecmd "github.com/kcp-dev/kcp/pkg/cliplugins/workspace/cmd"
	"github.com/kcp-dev/kcp/pkg/cmd/help"
)

func KubectlKcpCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "kcp",
		Short: "kubectl plugin for KCP",
		Long: help.Doc(`
			KCP is the easiest way to manage Kubernetes applications against one or
			more clusters, by giving you a personal control plane that schedules your
			workloads onto one or many clusters, and making it simple to pick up and
			move. Advanced use cases including spreading your apps across clusters for
			resiliency, scheduling batch workloads onto clusters with free capacity,
			and enabling collaboration for individual teams without having access to
			the underlying clusters.

			This command provides KCP specific sub-command for kubectl.
		`),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// setup klog
	fs := goflags.NewFlagSet("klog", goflags.PanicOnError)
	klog.InitFlags(fs)
	root.PersistentFlags().AddGoFlagSet(fs)

	if v := version.Get().String(); len(v) == 0 {
		root.Version = "<unknown>"
	} else {
		root.Version = v
	}

	workspaceCmd, err := workspacecmd.New(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	root.AddCommand(workspaceCmd)

	workloadCmd, err := workloadcmd.New(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	root.AddCommand(workloadCmd)

	crdCmd := crdcmd.New(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	root.AddCommand(crdCmd)

	return root
}

func GenerateDocs(filepath string) error {
	return doc.GenMarkdownTree(KubectlKcpCommand(), filepath)
}
