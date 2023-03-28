/*
Copyright 2023.

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

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dinov1alpha "github.com/roehrich-hpe/conditions-array-play/api/v1alpha1"
)

var _ = Describe("Bird Controller", func() {

	var (
		bird *dinov1alpha.Bird
	)

	BeforeEach(func() {
		birdId := uuid.NewString()[0:8]
		bird = &dinov1alpha.Bird{
			ObjectMeta: metav1.ObjectMeta{
				Name:      birdId,
				Namespace: corev1.NamespaceDefault,
			},
			Spec: dinov1alpha.BirdSpec{
				Foo: "",
			},
		}
	})

	It("creates a beak resource", func() {
		Expect(k8sClient.Create(context.TODO(), bird)).To(Succeed())
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(bird), bird)).To(Succeed())
			g.Expect(meta.IsStatusConditionTrue(bird.Status.Conditions, dinov1alpha.BirdConditionBeakResource)).To(BeTrue())
		}).Should(Succeed())
	})
})
