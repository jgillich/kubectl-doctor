package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var flags struct {
	*genericclioptions.ConfigFlags
}

func main() {
	var rootCmd = &cobra.Command{Use: "kubectl-doctor"}
	rootCmd.AddCommand(cmdTriage, cmdRules)

	flags.ConfigFlags = genericclioptions.NewConfigFlags(true)
	flags.AddFlags(rootCmd.PersistentFlags())
	flags.AddFlags(cmdTriage.PersistentFlags())
	flags.AddFlags(cmdRules.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
