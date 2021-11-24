/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"os"
	"strings"
	"time"

	"stash.appscode.dev/apimachinery/apis"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/docker"
	"stash.appscode.dev/stash/pkg/backup"
	"stash.appscode.dev/stash/pkg/scale"
	"stash.appscode.dev/stash/pkg/util"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/meta"
)

func NewCmdBackup() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string
		opt            = backup.Options{
			Namespace:      meta.Namespace(),
			ScratchDir:     "/tmp",
			PodLabelsPath:  "/etc/stash/labels",
			DockerRegistry: docker.ACRegistry,
			QPS:            100,
			Burst:          100,
			ResyncPeriod:   5 * time.Minute,
			MaxNumRequeues: 5,
			NumThreads:     1,
		}
	)

	cmd := &cobra.Command{
		Use:               "backup",
		Short:             "Run Stash Backup",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				klog.Fatalf("Could not get Kubernetes config: %s", err)
			}
			kubeClient := kubernetes.NewForConfigOrDie(config)
			stashClient := cs.NewForConfigOrDie(config)

			opt.NodeName = os.Getenv("NODE_NAME")
			if opt.NodeName == "" {
				klog.Fatalln(`Missing ENV var "NODE_NAME"`)
			}
			opt.PodName = os.Getenv("POD_NAME")
			if opt.PodName == "" {
				klog.Fatalln(`Missing ENV var "POD_NAME"`)
			}

			if err := opt.Workload.Canonicalize(); err != nil {
				klog.Fatalf(err.Error())
			}
			if opt.SnapshotHostname, opt.SmartPrefix, err = opt.Workload.HostnamePrefix(opt.PodName, opt.NodeName); err != nil {
				klog.Fatalf(err.Error())
			}
			if err = util.WorkloadExists(kubeClient, opt.Namespace, opt.Workload); err != nil {
				klog.Fatalf(err.Error())
			}
			opt.ScratchDir = strings.TrimSuffix(opt.ScratchDir, "/") // make ScratchDir in setup()

			ctrl := backup.New(kubeClient, stashClient, opt)

			if opt.RunViaCron {
				klog.Infoln("Running backup periodically via cron")
				if err = ctrl.BackupScheduler(); err != nil {
					klog.Fatal(err)
				}
			} else { // for offline backup
				if opt.Workload.Kind == apis.KindDaemonSet || opt.Workload.Kind == apis.KindStatefulSet {
					klog.Infoln("Running backup once")
					if err = ctrl.Backup(); err != nil {
						klog.Fatal(err)
					}
				} else {
					//if replica > 1 we should not take backup
					replica, err := util.WorkloadReplicas(kubeClient, opt.Namespace, opt.Workload.Kind, opt.Workload.Name)
					if err != nil {
						klog.Fatal(err)
					}

					if replica > 1 {
						klog.Infof("Skipping backup...\n" +
							"Reason: Backup type offline and replica > 1\n" +
							"Backup has taken by another replica or scheduled CronJob hasn't run yet.")
					} else if !util.HasOldReplicaAnnotation(kubeClient, opt.Namespace, opt.Workload) {
						klog.Infof("Skipping backup...\n" +
							"Reason: Backup type offline and workload does not have 'old-replica' annotation.\n" +
							"Backup will be taken at next scheduled time.")
					} else {
						klog.Infoln("Running backup once")
						if err = ctrl.Backup(); err != nil {
							klog.Fatal(err)
						}

						// offline backup done. now scale up replica to original replica number
						err = scale.ScaleUpWorkload(kubeClient, opt)
						if err != nil {
							klog.Fatal(err)
						}
					}
				}
			}
			klog.Infoln("Exiting Stash Backup")
		},
	}
	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.Workload.Kind, "workload-kind", opt.Workload.Kind, "Kind of workload where sidecar pod is added.")
	cmd.Flags().StringVar(&opt.Workload.Name, "workload-name", opt.Workload.Name, "Name of workload where sidecar pod is added.")
	cmd.Flags().StringVar(&opt.ResticName, "restic-name", opt.ResticName, "Name of the Restic used as configuration.")
	cmd.Flags().StringVar(&opt.ScratchDir, "scratch-dir", opt.ScratchDir, "Directory used to store temporary files. Use an `emptyDir` in Kubernetes.")
	cmd.Flags().StringVar(&opt.PushgatewayURL, "pushgateway-url", opt.PushgatewayURL, "URL of Prometheus pushgateway used to cache backup metrics")
	cmd.Flags().Float64Var(&opt.QPS, "qps", opt.QPS, "The maximum QPS to the master from this client")
	cmd.Flags().IntVar(&opt.Burst, "burst", opt.Burst, "The maximum burst for throttle")
	cmd.Flags().DurationVar(&opt.ResyncPeriod, "resync-period", opt.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
	cmd.Flags().BoolVar(&opt.RunViaCron, "run-via-cron", opt.RunViaCron, "Run backup periodically via cron.")
	cmd.Flags().StringVar(&opt.DockerRegistry, "docker-registry", opt.DockerRegistry, "Check job image registry.")
	cmd.Flags().StringVar(&opt.ImageTag, "image-tag", opt.ImageTag, "Check job image tag.")

	return cmd
}
