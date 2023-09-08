package main

import (
	"fmt"
	"os"

	"github.com/jgillich/kubectl-doctor/pkg/triage"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/printers"
)

var cmdRules = &cobra.Command{
	Use:   "rules",
	Short: "Display triage rules.",
	Run: func(c *cobra.Command, args []string) {
		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()

		fmt.Fprintf(w, "RULE\tDESCRIPTION\n")
		for _, tr := range triage.List {
			fmt.Fprintf(w, "%s\t%s\n", triage.Id(tr), tr.Description())
		}
	},
}
