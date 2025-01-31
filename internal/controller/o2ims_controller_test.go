/*
Copyright 2025.

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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ovranv1alpha1 "github.com/example/o2-ims-operator/api/v1alpha1"
)

var (
	server     *ghttp.Server
	statusCode int
	body       []byte
	path       string
	addr       string
)

var _ = Describe("O2ims Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}

		BeforeEach(func() {
			By("launching mock REST server")
			server = ghttp.NewServer()
			statusCode = 200
			path = "/version"
			body = []byte(`{ "version": "1.0.0" }`)
			addr = "http://" + server.Addr() + path
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", path),
					ghttp.RespondWithPtr(&statusCode, &body),
				))

			By("creating the custom resource for the Kind O2ims")
			o2ims := &ovranv1alpha1.O2ims{}
			err := k8sClient.Get(ctx, typeNamespacedName, o2ims)
			if err != nil && errors.IsNotFound(err) {
				resource := &ovranv1alpha1.O2ims{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: ovranv1alpha1.O2imsSpec{
						Url: addr,
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &ovranv1alpha1.O2ims{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance O2ims")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())

			By("Cleanup the REST server mock")
			server.Close()
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &O2imsReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
