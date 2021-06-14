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

package ncsc

import (
	"context"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	autohealerv1alpha1 "github.com/stakater/jarvis/api/v1alpha1"
)

// NodeConditionSetReconciler reconciles a NodeConditionSet object
type NodeConditionSetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=autohealer.stakater.com,resources=nodeconditionsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autohealer.stakater.com,resources=nodeconditionsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autohealer.stakater.com,resources=nodeconditionsets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConditionSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *NodeConditionSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	reqLogger := r.Log.WithValues("nodeconditionset", req.NamespacedName)
	reqLogger.Info("NodeConditionSet reconciliation started...")

	ncs := &autohealerv1alpha1.NodeConditionSet{}
	err = r.Get(ctx, req.NamespacedName, ncs)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, errors.Wrap(err, "failed to load NodeConditionSets data")
	}

	reqLogger.Info("Found nodeConditionSet", "NodeConditionSet name", ncs.Name)

	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeConditionSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autohealerv1alpha1.NodeConditionSet{}).
		Complete(r)
}

func GetConditionSetMap(ctx context.Context, client client.Client) (*autohealerv1alpha1.NodeConditionSetList, error) {
	var ncsList autohealerv1alpha1.NodeConditionSetList
	err := client.List(ctx, &ncsList)
	if err != nil {
		return nil, err
	}
	return &ncsList, nil
}
