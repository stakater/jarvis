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

package ntc

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	autohealerv1alpha1 "github.com/stakater/jarvis/api/v1alpha1"
	"github.com/stakater/jarvis/controllers/ncsc"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NtcReconciler reconciles a Ntc object
type NtcReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=autohealer.stakater.com,resources=ntcs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autohealer.stakater.com,resources=ntcs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autohealer.stakater.com,resources=ntcs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ntc object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NtcReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ntcLogger := r.Log.WithValues("ntc", req.NamespacedName)
	ntcLogger.Info("NodeTaintController reconciliation started...")
	// your logic here

	var nodeList v1.NodeList
	err := r.List(context.Background(), &nodeList)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, errors.Wrap(err, "failed to load Nodes data")
	}

	csList, err := ncsc.GetConditionSetMap(ctx, r.Client)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, errors.Wrap(err, "failed to load ConditionSets data")
	}

	for _, node := range nodeList.Items {
		ntcLogger.Info("Found node", "Node name", node.Status.Conditions)
		checkForNodeConditionMatchingConditionSet(csList, node.Status.Conditions)
	}

	return ctrl.Result{}, nil
}

func checkForNodeConditionMatchingConditionSet(ncsList *autohealerv1alpha1.NodeConditionSetList, nodeConditions []v1.NodeCondition) {

	for _, ncs := range ncsList.Items {
		_ = ncs.Spec.NodeConditions

	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *NtcReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&v1.Node{}).
		Owns(&autohealerv1alpha1.NodeConditionSet{}).
		Complete(r)
}
