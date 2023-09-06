package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/emirozer/kubectl-doctor/pkg/triage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	example = `
	# triage everything in the cluster
	kubectl doctor
`
	longDesc = `
    kubectl-doctor plugin will scan the given k8s cluster for any kind of anomalies and reports back to its user.
    example anomalies:
        * deployments that are older than 30d with 0 available,
        * deployments that do not have minimum availability,
        * kubernetes nodes cpu usage or memory usage too high. or too low to report scaledown possiblity
`
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)
}

// DoctorOptions specify what the doctor is going to do
type DoctorOptions struct {

	// Doctor options
	Flags   *genericclioptions.ConfigFlags
	Context context.Context
}

// NewDoctorOptions new doctor options initializer
func NewDoctorOptions() *DoctorOptions {
	return &DoctorOptions{
		Flags:   genericclioptions.NewConfigFlags(true),
		Context: context.Background(),
	}
}

// NewDoctorCmd returns a cobra command wrapping DoctorOptions
func NewDoctorCmd() *cobra.Command {

	opts := NewDoctorOptions()

	cmd := &cobra.Command{
		Use:     "doctor",
		Short:   "start triage for current targeted kubernetes cluster",
		Long:    longDesc,
		Example: example,
		Run: func(c *cobra.Command, args []string) {
			cmdutil.CheckErr(opts.Run(args))
		},
	}

	opts.Flags.AddFlags(cmd.Flags())

	return cmd
}

// Run doctor run
func (o *DoctorOptions) Run(args []string) error {
	// var kubeconfig *string
	// chain := clientcmd.NewDefaultClientConfigLoadingRules().Precedence
	// if len(chain) > 0 {
	// 	kubeconfig = &chain[0]
	// }
	// kubeconfig = flag.String(
	// 	"kubeconfig",
	// 	*kubeconfig,
	// 	"(optional) absolute path to the kubeconfig file",
	// )
	// var kubecontext = flag.String(
	// 	"context",
	// 	"",
	// 	"(optional) name of kube context",
	// )

	// flag.Parse()

	// log.Info(*kubeconfig)

	// config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
	// 	&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig},
	// 	&clientcmd.ConfigOverrides{
	// 		CurrentContext: *kubecontext,
	// 	}).ClientConfig()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	restConfig, err := o.Flags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return err
	}

	cl, err := client.New(restConfig, client.Options{})
	if err != nil {
		return err
	}
	for _, t := range triage.List {
		anomalies, err := t.Triage(context.TODO(), cl)
		if err != nil {
			log.Error(fmt.Errorf("%s: %w", t.Id(), err))
		} else {
			if len(anomalies) > 0 {
				fmt.Printf("%s: ", t.Id())
				for _, a := range anomalies {
					fmt.Printf("%s ", a.Name)
				}
				fmt.Println()
			}
		}
	}
	return nil
}
