package main

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/kubectl/pkg/cmd/util"
)

var flags struct {
	*genericclioptions.ConfigFlags
}

func main() {
	var rootCmd = &cobra.Command{Use: "kubectl-doctor"}
	rootCmd.AddCommand(cmdTriage)

	flags.ConfigFlags = genericclioptions.NewConfigFlags(true)
	flags.AddFlags(rootCmd.PersistentFlags())
	flags.AddFlags(cmdTriage.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		util.CheckErr(err)
	}
}
