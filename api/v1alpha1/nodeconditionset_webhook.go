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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var nodeconditionsetlog = logf.Log.WithName("nodeconditionset-resource")

func (ncs *NodeConditionSet) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(ncs).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-autohealer-stakater-com-v1alpha1-nodeconditionset,mutating=true,failurePolicy=fail,sideEffects=None,groups=autohealer.stakater.com,resources=nodeconditionsets,verbs=create;update,versions=v1alpha1,name=mnodeconditionset.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &NodeConditionSet{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (ncs *NodeConditionSet) Default() {
	nodeconditionsetlog.Info("default", "name", ncs.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-autohealer-stakater-com-v1alpha1-nodeconditionset,mutating=false,failurePolicy=fail,sideEffects=None,groups=autohealer.stakater.com,resources=nodeconditionsets,verbs=create;update,versions=v1alpha1,name=vnodeconditionset.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &NodeConditionSet{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (ncs *NodeConditionSet) ValidateCreate() error {
	nodeconditionsetlog.Info("validate create", "name", ncs.Name)
	status, err := ncs.Validate()
	if !status {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (ncs *NodeConditionSet) ValidateUpdate(old runtime.Object) error {
	nodeconditionsetlog.Info("validate update", "name", ncs.Name)
	status, err := ncs.Validate()
	if !status {
		return err
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (ncs *NodeConditionSet) ValidateDelete() error {
	nodeconditionsetlog.Info("validate delete", "name", ncs.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
