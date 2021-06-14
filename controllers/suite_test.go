/*
Copyright 2021.

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

package controllers

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	autohealerv1alpha1 "github.com/stakater/jarvis/api/v1alpha1"
	"github.com/stakater/jarvis/controllers/ncsc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	logf.Log.Info("bootstrapping test environment.............")
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = autohealerv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:    scheme.Scheme,
		Namespace: "",
	})

	err = (&ncsc.NodeConditionSetReconciler{
		Client: k8sClient,
		Log:    ctrl.Log.WithName("controllers").WithName("Ntc"),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	ncs := createNodeConditionSet()
	Expect(k8sClient.Create(context.Background(), ncs)).Should(Succeed())

	ncs1 := &autohealerv1alpha1.NodeConditionSet{}
	err = k8sClient.Get(context.Background(), types.NamespacedName{
		Namespace: "default",
		Name:      "nodeconditionset-1",
	}, ncs1)

	fmt.Println("NodeConditonSet.............", "NodeConditionSet", ncs)
	if err != nil {
		logf.Log.Error(err, "Failed to get ConditionSet by Name: conditionset-1")
	} else {
		logf.Log.Info("Found NodeCondtionSet, ", "CondtionSet: ", ncs1)
	}

}, 60)

func createNodeConditionSet() *autohealerv1alpha1.NodeConditionSet {
	ncs := &autohealerv1alpha1.NodeConditionSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeConditionSet",
			APIVersion: "autohealer.stakater.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nodeconditionset-1",
			Namespace: "default",
		},
		Spec: autohealerv1alpha1.NodeConditionSetSpec{
			Name:     "KernelDeadlock",
			Effect:   "NoExecute",
			TaintKey: "node.stakater.com/KernelDeadlock",
			NodeConditions: []autohealerv1alpha1.NodeCondition{
				{
					Type:   "KernelDeadlock",
					Status: "True",
				},
			},
		},
	}

	return ncs
}

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
