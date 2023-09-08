package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jgillich/kubectl-doctor/pkg/triage"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cmdTriage = &cobra.Command{
	Use:   "triage",
	Short: "Triage kubernetes cluster.",
	Run: func(c *cobra.Command, args []string) {
		restConfig, err := flags.ToRawKubeConfigLoader().ClientConfig()
		checkErr(err)

		cl, err := client.New(restConfig, client.Options{})
		checkErr(err)

		var report = map[triage.Triage][]triage.Anomaly{}

		for _, t := range triage.List {
			anomalies, err := t.Triage(c.Context(), cl)
			checkErr(err)
			report[t] = anomalies
		}

		w := printers.GetNewTabWriter(os.Stdout)
		defer w.Flush()
		fmt.Fprint(w, strings.Join([]string{"RULE", "SEVERITY", "NAMESPACE", "NAME", "REASON", "\n"}, "\t"))

		for tr, anomalies := range report {
			for _, anomaly := range anomalies {
				fmt.Fprint(w, strings.Join([]string{triage.Id(tr), tr.Severity().String(), anomaly.NamespacedName.Namespace, anomaly.NamespacedName.Name, anomaly.Reason, "\n"}, "\t"))
			}
		}
	},
}
