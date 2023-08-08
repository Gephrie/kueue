/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rayjob

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	testingjob "sigs.k8s.io/kueue/pkg/util/testingjobs/rayjob"
	"sigs.k8s.io/kueue/test/integration/framework"
	"sigs.k8s.io/kueue/test/util"
)

var _ = ginkgo.Describe("RayJob Webhook", func() {
	var ns *corev1.Namespace

	ginkgo.When("With manageJobsWithoutQueueName disabled", ginkgo.Ordered, ginkgo.ContinueOnFailure, func() {

		ginkgo.BeforeAll(func() {
			fwk = &framework.Framework{
				ManagerSetup: managerSetup(),
				CRDPath:      crdPath,
				DepCRDPaths:  []string{rayCrdPath},
				WebhookPath:  webhookPath,
			}
			ctx, cfg, k8sClient = fwk.Setup()
		})
		ginkgo.BeforeEach(func() {
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "rayjob-",
				},
			}
			gomega.Expect(k8sClient.Create(ctx, ns)).To(gomega.Succeed())
		})

		ginkgo.AfterEach(func() {
			gomega.Expect(util.DeleteNamespace(ctx, k8sClient, ns)).To(gomega.Succeed())
		})
		ginkgo.AfterAll(func() {
			fwk.Teardown()
		})

		ginkgo.It("the creation doesn't succeed if the queue name is invalid", func() {
			job := testingjob.MakeJob(jobName, ns.Name).Queue("indexed_job").Obj()
			err := k8sClient.Create(ctx, job)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(apierrors.IsForbidden(err)).Should(gomega.BeTrue(), "error: %v", err)
		})

		ginkgo.It("invalid configuration shutdown after job finishes", func() {
			job := testingjob.MakeJob(jobName, ns.Name).
				Queue("queue-name").
				ShutdownAfterJobFinishes(false).
				Obj()
			err := k8sClient.Create(ctx, job)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(apierrors.IsForbidden(err)).Should(gomega.BeTrue(), "error: %v", err)
		})
	})
})
