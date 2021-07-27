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

package attractor

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	meridiov1alpha1 "github.com/nordix/meridio-operator/api/v1alpha1"
	common "github.com/nordix/meridio-operator/controllers/common"
)

// AttractorReconciler reconciles a Attractor object
type AttractorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=meridio.nordix.org,resources=attractors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=meridio.nordix.org,resources=attractors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=meridio.nordix.org,resources=attractors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *AttractorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	attr := &meridiov1alpha1.Attractor{}
	executor := common.NewExecutor(r.Scheme, r.Client, ctx, nil, r.Log)

	err := r.Get(ctx, req.NamespacedName, attr)
	if err != nil {
		if apierrors.IsNotFound(err) {
			if err != nil {
				return ctrl.Result{}, nil
			}
			return reconcile.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	attr.Status = meridiov1alpha1.AttractorStatus{}

	trench, err := validateAttractor(executor, attr)
	if err != nil {
		return ctrl.Result{}, err
	}
	// create/update load-balancer & nse-vlan deployment
	executor.SetOwner(attr)
	lb, err := NewLoadBalancer(executor, attr, trench)
	if err != nil {
		return ctrl.Result{}, err
	}
	nse, err := NewNSE(executor, attr)
	if err != nil {
		return ctrl.Result{}, err
	}
	alb, err := lb.getAction(executor)
	if err != nil {
		return ctrl.Result{}, err
	}
	anse, err := nse.getAction(executor)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = executor.RunAll(append(alb, anse...))
	if err != nil {
		return ctrl.Result{}, err
	}
	// update attractor
	executor.SetOwner(trench)
	actions := getAttractorActions(executor, attr)
	err = executor.RunAll(actions)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func validateAttractor(e *common.Executor, attr *meridiov1alpha1.Attractor) (*meridiov1alpha1.Trench, error) {
	// get the trench by attractor label
	selector := client.ObjectKey{
		Namespace: attr.ObjectMeta.Namespace,
		Name:      attr.ObjectMeta.Labels["trench"],
	}
	trench, err := common.GetTrenchbySelector(e, selector)
	if err != nil {
		if apierrors.IsNotFound(err) {
			msg := "labeled trench not found"
			attr.Status.LB = meridiov1alpha1.LBStatus.Disengaged
			attr.Status.Message = msg
			return nil, nil
		} else {
			return nil, err
		}
	}
	// validation: get the all attractors with same trench, verdict the number should not be greater than 1
	al := &meridiov1alpha1.AttractorList{}
	sel := labels.Set{"trench": trench.ObjectMeta.Name}
	err = e.ListObject(al, &client.ListOptions{
		LabelSelector: sel.AsSelector(),
		Namespace:     attr.ObjectMeta.Namespace,
	})
	if err != nil {
		msg := "at least one attractor should be found"
		attr.Status.Status = meridiov1alpha1.BaseStatus.Error
		attr.Status.Message = msg
	}
	if len(al.Items) > 1 {
		msg := "only one attractor is supported per trench"
		attr.Status.Status = meridiov1alpha1.ConfigStatus.Rejected
		attr.Status.Message = msg
	}
	if attr.Status.Status == meridiov1alpha1.BaseStatus.NoPhase {
		attr.Status.Status = meridiov1alpha1.ConfigStatus.Accepted
		if attr.Status.LB == meridiov1alpha1.BaseStatus.NoPhase {
			attr.Status.LB = meridiov1alpha1.LBStatus.Engaged
		}
	}
	return trench, nil
}

func getAttractorActions(e *common.Executor, attr *meridiov1alpha1.Attractor) []common.Action {
	var actions []common.Action
	// set the status for the vip
	attrnsname := fmt.Sprintf("%s/%s", attr.GetNamespace(), attr.GetName())
	actions = append(actions, common.NewUpdateStatusAction(attr, fmt.Sprintf("update %s status: %v", attrnsname, attr.Status.Status)))
	// if attr doesn't have an existing trench, update the status only
	if e.GetOwner().(*meridiov1alpha1.Trench) == nil {
		return actions
	}
	// if labeled trench exsits, update the ownerReference
	actions = append(actions, common.NewUpdateAction(attr, fmt.Sprintf("update %s ownerReference: %s", attrnsname, e.GetOwner().GetName())))
	return actions
}

// SetupWithManager sets up the controller with the Manager.
func (r *AttractorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&meridiov1alpha1.Attractor{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
