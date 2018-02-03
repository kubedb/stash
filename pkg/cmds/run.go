package cmds

import (
	"fmt"
	"net/http"
	"time"

	"github.com/appscode/go/log"
	stringz "github.com/appscode/go/strings"
	v "github.com/appscode/go/version"
	"github.com/appscode/kutil/discovery"
	"github.com/appscode/pat"
	api "github.com/appscode/stash/apis/stash"
	cs "github.com/appscode/stash/client"
	"github.com/appscode/stash/pkg/controller"
	"github.com/appscode/stash/pkg/docker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdRun() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string
		address        = ":56790"
		opts           = controller.Options{
			DockerRegistry: docker.ACRegistry,
			StashImageTag:  stringz.Val(v.Version.Version, "canary"),
			ResyncPeriod:   10 * time.Minute,
			MaxNumRequeues: 5,
			NumThreads:     2,
		}
		scratchDir = "/tmp"
	)

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Run Stash operator",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				log.Fatalln(err)
			}
			kubeClient := kubernetes.NewForConfigOrDie(config)
			stashClient := cs.NewForConfigOrDie(config)
			crdClient := crd_cs.NewForConfigOrDie(config)

			// get kube api server version
			opts.KubectlImageTag, err = discovery.GetBaseVersion(kubeClient.Discovery())
			if err != nil {
				log.Fatalf("Failed to detect server version, reason: %s\n", err)
			}

			ctrl := controller.New(kubeClient, crdClient, stashClient, opts)
			err = ctrl.Setup()
			if err != nil {
				log.Fatalln(err)
			}

			log.Infof("Starting operator version %s+%s ...", v.Version.Version, v.Version.CommitHash)
			// Now let's start the controller
			stop := make(chan struct{})
			defer close(stop)
			go ctrl.Run(1, stop)

			m := pat.New()
			m.Get("/metrics", promhttp.Handler())

			pattern := fmt.Sprintf("/%s/v1beta1/namespaces/%s/restics/%s/metrics", api.GroupName, PathParamNamespace, PathParamName)
			log.Infof("URL pattern: %s", pattern)
			exporter := &PrometheusExporter{
				kubeClient:  kubeClient,
				stashClient: stashClient.StashV1alpha1(),
				scratchDir:  scratchDir,
			}
			m.Get(pattern, exporter)

			http.Handle("/", m)
			log.Infoln("Listening on", address)
			log.Fatal(http.ListenAndServe(address, nil))
		},
	}
	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&address, "address", address, "Address to listen on for web interface and telemetry.")
	cmd.Flags().BoolVar(&opts.EnableRBAC, "rbac", opts.EnableRBAC, "Enable RBAC for operator")
	cmd.Flags().StringVar(&scratchDir, "scratch-dir", scratchDir, "Directory used to store temporary files. Use an `emptyDir` in Kubernetes.")
	cmd.Flags().DurationVar(&opts.ResyncPeriod, "resync-period", opts.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
	cmd.Flags().StringVar(&opts.StashImageTag, "image-tag", opts.StashImageTag, "Image tag for sidecar, init-container, check-job and recovery-job")
	cmd.Flags().StringVar(&opts.DockerRegistry, "docker-registry", opts.DockerRegistry, "Docker image registry for sidecar, init-container, check-job, recovery-job and kubectl-job")

	return cmd
}
