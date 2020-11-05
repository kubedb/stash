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
	"fmt"
	"strings"
	"time"

	v1beta1_api "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/invoker"
	"stash.appscode.dev/apimachinery/pkg/restic"
	"stash.appscode.dev/stash/pkg/restore"
	"stash.appscode.dev/stash/pkg/util"

	"github.com/spf13/cobra"
	"gomodules.xyz/x/log"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kmodules.xyz/client-go/meta"
)

func NewCmdRestore() *cobra.Command {
	opt := &restore.Options{
		MasterURL:      "",
		KubeconfigPath: "",
		Namespace:      meta.Namespace(),
		SetupOpt: restic.SetupOptions{
			ScratchDir:  "/tmp",
			EnableCache: true,
		},
		RestoreModel: restore.RestoreModelInitContainer,
	}

	cmd := &cobra.Command{
		Use:               "restore",
		Short:             "Restore from backup",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// create client
			config, err := clientcmd.BuildConfigFromFlags(opt.MasterURL, opt.KubeconfigPath)
			if err != nil {
				log.Fatal(err)
				return err
			}
			opt.Config = config
			opt.KubeClient = kubernetes.NewForConfigOrDie(config)
			opt.StashClient = cs.NewForConfigOrDie(config)
			opt.Metrics.JobName = fmt.Sprintf("%s-%s-%s", strings.ToLower(opt.InvokerKind), opt.Namespace, opt.InvokerName)

			inv, err := invoker.ExtractRestoreInvokerInfo(opt.KubeClient, opt.StashClient, opt.InvokerKind, opt.InvokerName, opt.Namespace)
			if err != nil {
				return err
			}

			for _, targetInfo := range inv.TargetsInfo {
				if targetInfo.Target != nil && targetMatched(targetInfo.Target.Ref, opt.TargetKind, opt.TargetName) {

					// Ensure restore order
					if inv.ExecutionOrder == v1beta1_api.Sequential {
						err = waitUntilAllPreviousTargetsExecuted(opt, targetInfo.Target.Ref)
						if err != nil {
							return err
						}
					}

					opt.Host, err = util.GetHostName(targetInfo.Target)
					if err != nil {
						return err
					}

					// run restore
					restoreOutput, restoreErr := opt.Restore(inv, targetInfo)
					if restoreErr != nil {
						err = opt.HandleRestoreFailure(inv, targetInfo, restoreErr)
						return errors.NewAggregate([]error{restoreErr, err})

					}
					if restoreOutput != nil {
						return opt.HandleRestoreSuccess(restoreOutput, inv, targetInfo)
					}
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&opt.MasterURL, "master", opt.MasterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&opt.KubeconfigPath, "kubeconfig", opt.KubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.InvokerKind, "invoker-kind", opt.InvokerKind, "Kind of the respective restore invoker")
	cmd.Flags().StringVar(&opt.InvokerName, "invoker-name", opt.InvokerName, "Name of the respective restore invoker")
	cmd.Flags().StringVar(&opt.TargetName, "target-name", opt.TargetName, "Name of the Target")
	cmd.Flags().StringVar(&opt.TargetKind, "target-kind", opt.TargetKind, "Kind of the Target")
	cmd.Flags().DurationVar(&opt.BackoffMaxWait, "backoff-max-wait", 0, "Maximum wait for initial response from kube apiserver; 0 disables the timeout")
	cmd.Flags().BoolVar(&opt.SetupOpt.EnableCache, "enable-cache", opt.SetupOpt.EnableCache, "Specify whether to enable caching for restic")
	cmd.Flags().Int64Var(&opt.SetupOpt.MaxConnections, "max-connections", opt.SetupOpt.MaxConnections, "Specify maximum concurrent connections for GCS, Azure and B2 backend")
	cmd.Flags().StringVar(&opt.SetupOpt.SecretDir, "secret-dir", opt.SetupOpt.SecretDir, "Directory where storage secret has been mounted")

	cmd.Flags().BoolVar(&opt.Metrics.Enabled, "metrics-enabled", opt.Metrics.Enabled, "Specify whether to export Prometheus metrics")
	cmd.Flags().StringVar(&opt.Metrics.PushgatewayURL, "pushgateway-url", opt.Metrics.PushgatewayURL, "Pushgateway URL where the metrics will be pushed")
	cmd.Flags().StringVar(&opt.RestoreModel, "restore-model", opt.RestoreModel, "Specify whether using job or init-container to restore (default init-container)")

	return cmd
}

func waitUntilAllPreviousTargetsExecuted(opt *restore.Options, tref v1beta1_api.TargetRef) error {
	return wait.PollImmediate(5*time.Second, 30*time.Minute, func() (bool, error) {
		log.Infof("Waiting for all previous targets to complete their restore process...")
		inv, err := invoker.ExtractRestoreInvokerInfo(opt.KubeClient, opt.StashClient, opt.InvokerKind, opt.InvokerName, opt.Namespace)
		if err != nil {
			return false, err
		}
		return inv.NextInOrder(tref, inv.Status.TargetStatus), nil
	})
}
