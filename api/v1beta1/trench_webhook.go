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

package v1beta1

import (
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var trenchlog = logf.Log.WithName("trench-resource")

func (r *Trench) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-meridio-nordix-org-v1alpha1-trench,mutating=true,failurePolicy=fail,sideEffects=None,groups=meridio.nordix.org,resources=trenches,verbs=create;update,versions=v1alpha1,name=mtrench.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Trench{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Trench) Default() {
	trenchlog.Info("default", "name", r.Name)
	if r.Spec.Family == "" {
		r.Spec.Family = "dualstack"
	} else {
		r.Spec.Family = strings.ToLower(r.Spec.Family)
	}
}

//+kubebuilder:webhook:path=/validate-meridio-nordix-org-v1alpha1-trench,mutating=false,failurePolicy=fail,sideEffects=None,groups=meridio.nordix.org,resources=trenches,verbs=create;update,versions=v1alpha1,name=vtrench.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Trench{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Trench) ValidateCreate() error {
	trenchlog.Info("validate create", "name", r.Name)
	return r.validateTrench()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Trench) ValidateUpdate(old runtime.Object) error {
	trenchlog.Info("validate update", "name", r.Name)
	err := r.validateUpdate(old)
	if err != nil {
		return err
	}
	return r.validateTrench()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Trench) ValidateDelete() error {
	trenchlog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *Trench) validateTrench() error {
	var allErrs field.ErrorList

	if err := r.validateSpec(); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(r.GroupKind(), r.Name, allErrs)
}

func (r *Trench) validateSpec() field.ErrorList {
	var allErrs field.ErrorList
	if !IPFamily(r.Spec.Family).IsValid() {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("family"), r.Spec.Family, "invalid family"))
	}
	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func (r *Trench) validateUpdate(old runtime.Object) error {
	typedOld, ok := old.(*Trench)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a %s got a %T", r.GroupVersionKind().Kind, typedOld))
	}

	if r.Spec.Family != typedOld.Spec.Family {
		return apierrors.NewForbidden(r.GroupResource(),
			r.Name, field.Forbidden(field.NewPath("spec", "family"), "update on family is forbidden"))
	}
	return nil
}

// +kubebuilder:validation:Enum=ipv4;ipv6;dualstack

// IPFamily describes the traffic type in the trench
// Only one of the following ip family can be specified.
// If the traffic is IPv4 only, use IPv4, similarly,
// use IPv6 if the traffic is IPv6 only, otherwise, use
// dualstack which handles both IPv4 and IPv6 traffic.
type IPFamily string

const (
	IPv4      IPFamily = "ipv4"
	IPv6      IPFamily = "ipv6"
	Dualstack IPFamily = "dualstack"
)

// IsValid returns true if the receiver is a valid ipFamily type
func (f IPFamily) IsValid() bool {
	switch f {
	case IPv4, IPv6, Dualstack:
		return true
	default:
		return false
	}
}
