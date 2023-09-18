package main

import (
	"fmt"
	"os"
	"strings"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/jgillich/kubectl-doctor/pkg/triage"
	"github.com/spf13/cobra"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cmdTriage = &cobra.Command{
	Use:   "triage",
	Short: "Triage kubernetes cluster.",
	Run: func(c *cobra.Command, args []string) {
		restConfig, err := flags.ToRawKubeConfigLoader().ClientConfig()
		checkErr(err)

		checkErr(apiextensionsv1.AddToScheme(scheme.Scheme))
		checkErr(cnpgv1.AddToScheme(scheme.Scheme))
		checkErr(helmv2beta1.AddToScheme(scheme.Scheme))
		checkErr(kustomizev1.AddToScheme(scheme.Scheme))

		cl, err := client.New(restConfig, client.Options{
			Scheme: scheme.Scheme,
		})
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
