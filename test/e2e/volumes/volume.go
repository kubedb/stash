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

package volumes

import (
	"context"
	"fmt"

	"stash.appscode.dev/apimachinery/apis"
	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	"stash.appscode.dev/stash/test/e2e/framework"
	. "stash.appscode.dev/stash/test/e2e/matcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Volume", func() {
	var f *framework.Invocation

	BeforeEach(func() {
		f = framework.NewInvocation()
	})

	JustAfterEach(func() {
		f.PrintDebugInfoOnFailure()
	})

	AfterEach(func() {
		err := f.CleanupTestResources()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("PVC", func() {
		Context("Restore in same PVC", func() {
			It("should Backup & Restore in the source PVC", func() {
				// Create new PVC
				pvc, err := f.CreateNewPVC(fmt.Sprintf("%s-%s", framework.SourceVolume, f.App()))
				Expect(err).NotTo(HaveOccurred())

				// Deploy a Pod
				pod, err := f.DeployPod(pvc.Name)
				Expect(err).NotTo(HaveOccurred())

				// Generate Sample Data
				sampleData, err := f.GenerateSampleData(pod.ObjectMeta, apis.KindPod)
				Expect(err).NotTo(HaveOccurred())

				// Setup a Minio Repository
				repo, err := f.SetupMinioRepository()
				Expect(err).NotTo(HaveOccurred())
				f.AppendToCleanupList(repo)

				// Setup PVC Backup
				backupConfig, err := f.SetupPVCBackup(pvc, repo)
				Expect(err).NotTo(HaveOccurred())

				// Take an Instant Backup of the Sample Data
				backupSession, err := f.TakeInstantBackup(backupConfig.ObjectMeta, v1beta1.BackupInvokerRef{
					Name: backupConfig.Name,
					Kind: v1beta1.ResourceKindBackupConfiguration,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that BackupSession has succeeded")
				completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionSucceeded))

				// Simulate disaster scenario. Delete the data from source PVC
				By("Deleting sample data from source Pod")
				err = f.CleanupSampleDataFromWorkload(pod.ObjectMeta, apis.KindPod)
				Expect(err).NotTo(HaveOccurred())

				// Restore the backed up data
				By("Restoring the backed up data")
				restoreSession, err := f.SetupRestoreProcessForPVC(pvc, repo)
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that RestoreSession succeeded")
				completedRS, err := f.StashClient.StashV1beta1().RestoreSessions(restoreSession.Namespace).Get(context.TODO(), restoreSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedRS.Status.Phase).Should(Equal(v1beta1.RestoreSucceeded))

				// Get restored data
				restoredData := f.RestoredData(pod.ObjectMeta, apis.KindPod)

				// Verify that restored data is same as the original data
				By("Verifying restored data is same as the original data")
				Expect(restoredData).Should(BeSameAs(sampleData))
			})
		})

		Context("Restore in different PVC", func() {
			It("should restore backed up data into different PVC", func() {
				// Create new PVC
				pvc, err := f.CreateNewPVC(fmt.Sprintf("%s-%s", framework.SourceVolume, f.App()))
				Expect(err).NotTo(HaveOccurred())

				// Deploy a Pod
				pod, err := f.DeployPod(pvc.Name)
				Expect(err).NotTo(HaveOccurred())

				// Generate Sample Data
				sampleData, err := f.GenerateSampleData(pod.ObjectMeta, apis.KindPod)
				Expect(err).NotTo(HaveOccurred())

				// Setup a Minio Repository
				repo, err := f.SetupMinioRepository()
				Expect(err).NotTo(HaveOccurred())
				f.AppendToCleanupList(repo)

				// Setup PVC Backup
				backupConfig, err := f.SetupPVCBackup(pvc, repo)
				Expect(err).NotTo(HaveOccurred())

				// Take an Instant Backup of the Sample Data
				backupSession, err := f.TakeInstantBackup(backupConfig.ObjectMeta, v1beta1.BackupInvokerRef{
					Name: backupConfig.Name,
					Kind: v1beta1.ResourceKindBackupConfiguration,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that BackupSession has succeeded")
				completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionSucceeded))

				// Create restored Pvc
				restoredPVC, err := f.CreateNewPVC(fmt.Sprintf("%s-%s", framework.RestoredVolume, f.App()))
				Expect(err).NotTo(HaveOccurred())

				// Deploy another Pod
				restoredPod, err := f.DeployPod(restoredPVC.Name)
				Expect(err).NotTo(HaveOccurred())

				// Restore the backed up data
				By("Restoring the backed up data")
				restoreSession, err := f.SetupRestoreProcessForPVC(restoredPVC, repo)
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that RestoreSession succeeded")
				completedRS, err := f.StashClient.StashV1beta1().RestoreSessions(restoreSession.Namespace).Get(context.TODO(), restoreSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedRS.Status.Phase).Should(Equal(v1beta1.RestoreSucceeded))

				// Get restored data
				restoredData := f.RestoredData(restoredPod.ObjectMeta, apis.KindPod)

				// Verify that restored data is same as the original data
				By("Verifying restored data is same as the original data")
				Expect(restoredData).Should(BeSameAs(sampleData))
			})
		})
	})
})
