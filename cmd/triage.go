package main

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jgillich/kubectl-doctor/pkg/triage"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cmdTriage = &cobra.Command{
	Use:   "triage",
	Short: "Triage kubernetes cluster.",
	Run: func(c *cobra.Command, args []string) {
		restConfig, err := flags.ToRawKubeConfigLoader().ClientConfig()
		util.CheckErr(err)

		cl, err := client.New(restConfig, client.Options{})
		util.CheckErr(err)

		var report = map[triage.Triage][]triage.Anomaly{}

		for _, t := range triage.List {
			anomalies, err := t.Triage(c.Context(), cl)
			util.CheckErr(err)
			report[t] = anomalies
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.Style{Format: table.FormatOptionsDefault, Box: table.BoxStyle{PaddingLeft: " ", PaddingRight: " "}})
		t.AppendHeader(table.Row{"Type", "Severity", "Namespace", "Name", "Reason"})
		for triage, anomalies := range report {
			for _, anomaly := range anomalies {
				t.AppendRow(table.Row{reflect.TypeOf(triage).Elem().Name(), triage.Severity(), anomaly.NamespacedName.Namespace, anomaly.NamespacedName.Name, anomaly.Reason})
			}
		}
		t.Render()
	},
}
