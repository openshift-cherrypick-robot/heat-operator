/*
Copyright 2022.

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

//
// Generated by:
//
// operator-sdk create webhook --group heat --version v1beta1 --kind Heat --programmatic-validation --defaulting
//

package v1beta1

import (
	"fmt"

	"github.com/openstack-k8s-operators/lib-common/modules/common/service"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// HeatDefaults -
type HeatDefaults struct {
	APIContainerImageURL    string
	CfnAPIContainerImageURL string
	EngineContainerImageURL string
}

var heatDefaults HeatDefaults

// log is for logging in this package.
var heatlog = logf.Log.WithName("heat-resource")

// SetupHeatDefaults - initialize Heat spec defaults for use with either internal or external webhooks
func SetupHeatDefaults(defaults HeatDefaults) {
	heatDefaults = defaults
	heatlog.Info("Heat defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *Heat) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-heat-openstack-org-v1beta1-heat,mutating=true,failurePolicy=fail,sideEffects=None,groups=heat.openstack.org,resources=heats,verbs=create;update,versions=v1beta1,name=mheat.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Heat{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Heat) Default() {
	heatlog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this Heat spec
func (spec *HeatSpec) Default() {
	if spec.HeatAPI.ContainerImage == "" {
		spec.HeatAPI.ContainerImage = heatDefaults.APIContainerImageURL
	}
	if spec.HeatCfnAPI.ContainerImage == "" {
		spec.HeatCfnAPI.ContainerImage = heatDefaults.CfnAPIContainerImageURL
	}
	if spec.HeatEngine.ContainerImage == "" {
		spec.HeatEngine.ContainerImage = heatDefaults.EngineContainerImageURL
	}
}

// Default - set defaults for this Heat spec core. This version is called
// by the OpenStackControlplane
func (spec *HeatSpecCore) Default() {
	// nothing here yet
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-heat-openstack-org-v1beta1-heat,mutating=false,failurePolicy=fail,sideEffects=None,groups=heat.openstack.org,resources=heats,verbs=create;update,versions=v1beta1,name=vheat.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Heat{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Heat) ValidateCreate() (admission.Warnings, error) {
	heatlog.Info("validate create", "name", r.Name)

	var allErrs field.ErrorList
	basePath := field.NewPath("spec")
	if err := r.Spec.ValidateCreate(basePath); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "heat.openstack.org", Kind: "Heat"},
			r.Name, allErrs)
	}

	return nil, nil
}

// ValidateCreate - Exported function wrapping non-exported validate functions,
// this function can be called externally to validate an heat spec.
func (r *HeatSpec) ValidateCreate(basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatAPI").Child("override").Child("service"),
		r.HeatAPI.Override.Service)...)

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatCfnAPI").Child("override").Child("service"),
		r.HeatCfnAPI.Override.Service)...)

	return allErrs
}

func (r *HeatSpecCore) ValidateCreate(basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatAPI").Child("override").Child("service"),
		r.HeatAPI.Override.Service)...)

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatCfnAPI").Child("override").Child("service"),
		r.HeatCfnAPI.Override.Service)...)

	return allErrs
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Heat) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	heatlog.Info("validate update", "name", r.Name)
	oldHeat, ok := old.(*Heat)
	if !ok {
		return nil, apierrors.NewInternalError(
			fmt.Errorf("Expected a Heatv1 object, but got %T", oldHeat))
	}

	var allErrs field.ErrorList
	basePath := field.NewPath("spec")

	annotations := r.GetAnnotations()
	if err := r.Spec.ValidateUpdate(oldHeat.Spec, basePath, annotations); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "heat.openstack.org", Kind: "Heat"},
			r.Name, allErrs)
	}

	return nil, nil
}

// ValidateUpdate - Exported function wrapping non-exported validate functions,
// this function can be called externally to validate an barbican spec.
func (r *HeatSpec) ValidateUpdate(old HeatSpec, basePath *field.Path, annotations map[string]string) field.ErrorList {
	var allErrs field.ErrorList

	// Allow users to bypass this validation in cases where they have independently verified
	// the validity of their new database to ensure consistency with the current one.
	if _, ok := annotations[HeatDatabaseMigrationAnnotation]; !ok {
		// We currently have no logic in place to perform database migrations. Changing databases
		// would render all of the existing stacks unmanageable. We should block changes to the
		// databaseInstance to protect existing workloads.
		if old.DatabaseInstance != "" && r.DatabaseInstance != old.DatabaseInstance {
			allErrs = append(allErrs, field.Forbidden(
				field.NewPath("spec.databaseInstance"),
				"Changing the DatabaseInstance is not supported for existing deployments"))
		}
	}

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatAPI").Child("override").Child("service"),
		r.HeatAPI.Override.Service)...)

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatCfnAPI").Child("override").Child("service"),
		r.HeatCfnAPI.Override.Service)...)

	return allErrs
}

func (r *HeatSpecCore) ValidateUpdate(old HeatSpecCore, basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// We currently have no logic in place to perform database migrations. Changing databases
	// would render all of the existing stacks unmanageable. We should block changes to the
	// databaseInstance to protect existing workloads.
	if old.DatabaseInstance != "" && r.DatabaseInstance != old.DatabaseInstance {
		allErrs = append(allErrs, field.Forbidden(
			field.NewPath("spec.databaseInstance"),
			"Changing the DatabaseInstance is not supported for existing deployments"))
	}

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatAPI").Child("override").Child("service"),
		r.HeatAPI.Override.Service)...)

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("heatCfnAPI").Child("override").Child("service"),
		r.HeatCfnAPI.Override.Service)...)

	return allErrs
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Heat) ValidateDelete() (admission.Warnings, error) {
	heatlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

// SetDefaultRouteAnnotations sets HAProxy timeout values of the route
// NOTE: it is used by ctlplane webhook on openstack-operator
func (spec *HeatSpecCore) SetDefaultRouteAnnotations(annotations map[string]string) {
	const haProxyAnno = "haproxy.router.openshift.io/timeout"
	// Use a custom annotation to flag when the operator has set the default HAProxy timeout
	// With the annotation func determines when to overwrite existing HAProxy timeout with the APITimeout
	const heatAnno = "api.heat.openstack.org/timeout"

	valHeat, okHeat := annotations[heatAnno]
	valHAProxy, okHAProxy := annotations[haProxyAnno]

	// Human operator set the HAProxy timeout manually
	if !okHeat && okHAProxy {
		return
	}

	// Human operator modified the HAProxy timeout manually without removing the Heat flag
	if okHeat && okHAProxy && valHeat != valHAProxy {
		delete(annotations, heatAnno)
		return
	}

	timeout := fmt.Sprintf("%ds", spec.APITimeout)
	annotations[heatAnno] = timeout
	annotations[haProxyAnno] = timeout
}
