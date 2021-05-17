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
var conditionsetlog = logf.Log.WithName("conditionset-resource")

func (cs *ConditionSet) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(cs).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-autohealer-stakater-com-v1alpha1-conditionset,mutating=true,failurePolicy=fail,sideEffects=None,groups=autohealer.stakater.com,resources=conditionsets,verbs=create;update,versions=v1alpha1,name=mconditionset.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &ConditionSet{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (cs *ConditionSet) Default() {
	conditionsetlog.Info("default", "name", cs.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-autohealer-stakater-com-v1alpha1-conditionset,mutating=false,failurePolicy=fail,sideEffects=None,groups=autohealer.stakater.com,resources=conditionsets,verbs=create;update,versions=v1alpha1,name=vconditionset.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &ConditionSet{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (cs *ConditionSet) ValidateCreate() error {
	conditionsetlog.Info("validate create", "name", cs.Name)
	status, err := cs.Validate()
	if !status {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (cs *ConditionSet) ValidateUpdate(old runtime.Object) error {
	conditionsetlog.Info("validate update", "name", cs.Name)
	status, err := cs.Validate()
	if !status {
		return err
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (cs *ConditionSet) ValidateDelete() error {
	conditionsetlog.Info("validate delete", "name", cs.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
