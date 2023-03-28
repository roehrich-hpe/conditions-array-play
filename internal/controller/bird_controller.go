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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dinov1alpha "github.com/roehrich-hpe/conditions-array-play/api/v1alpha1"
)

// BirdReconciler reconciles a Bird object
type BirdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const finalizer = "bird.dino.hpe.com"

//+kubebuilder:rbac:groups=dino.hpe.com,resources=birds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dino.hpe.com,resources=birds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dino.hpe.com,resources=birds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bird object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *BirdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	bird := &dinov1alpha.Bird{}
	if err := r.Get(ctx, req.NamespacedName, bird); err != nil {
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !bird.GetDeletionTimestamp().IsZero() {
		log.Info("deleting resource")

		if controllerutil.ContainsFinalizer(bird, finalizer) {
			controllerutil.RemoveFinalizer(bird, finalizer)

			if err := r.Update(ctx, bird); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(bird, finalizer) {
		controllerutil.AddFinalizer(bird, finalizer)
		if err := r.Update(ctx, bird); err != nil {
			return ctrl.Result{}, err
		}

		// An update here will cause the reconciler to run again after k8s has
		// recorded the resource in its database.
		return ctrl.Result{}, nil
	}

	// If we don't have a Beak resource, do that now.  It'll share our name
	// and namespace.
	doUpdate, err := r.verifyBeakResource(ctx, bird)
	if err != nil {
		// return err=nil until we're ready to handle looping
		return ctrl.Result{}, nil
	}
	if doUpdate {
		if err := r.Status().Update(ctx, bird); err != nil {
			if !apierrors.IsConflict(err) {
				log.Error(err, "unable to update status")
				// return err=nil until we're ready to handle looping
				return ctrl.Result{}, nil
			}
			// Conflicts are a normal part of a reconciler's life. Try again.
			// This is common, so don't put this in the log.
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// verifyBeakResource will create a matching Beak resource, or update an
// existing one.
func (r *BirdReconciler) verifyBeakResource(ctx context.Context, bird *dinov1alpha.Bird) (bool, error) {
	log := log.FromContext(ctx)

	beak := &dinov1alpha.Beak{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bird.GetName(),
			Namespace: bird.GetNamespace(),
		},
		Spec: dinov1alpha.BeakSpec{
			Foo: "peck",
		},
	}
	result, err := ctrl.CreateOrUpdate(ctx, r.Client, beak,
		func() error {
			// Hook up garbage collection.
			return ctrl.SetControllerReference(bird, beak, r.Scheme)
		})

	if err != nil {
		log.Error(err, "unable to CreateOrUpdate the matching Beak resource")
		return false, err
	}

	if result == controllerutil.OperationResultCreated {
		log.Info("created matching Beak resource")
	} else if result == controllerutil.OperationResultNone {
		// no change
	} else {
		log.Info("updated matching Beak resource")
	}

	// The bird Condition array should reflect that we have this Beak
	// resource now.
	doUpdate := r.setBeakResourceCondition(bird)

	return doUpdate, nil
}

// setBeakResourceCondition will update the bird condition array to
// reflect that we have created the matching beak resource
func (r *BirdReconciler) setBeakResourceCondition(bird *dinov1alpha.Bird) bool {
	wantUpdate := false
	if !meta.IsStatusConditionTrue(bird.Status.Conditions, dinov1alpha.BirdConditionBeakResource) {
		cond := meta.FindStatusCondition(bird.Status.Conditions, dinov1alpha.BirdConditionBeakResource)
		if cond == nil {
			cond = &metav1.Condition{
				Type:    dinov1alpha.BirdConditionBeakResource,
				Message: "",
			}
		}
		cond.Reason = dinov1alpha.BirdConditionResourceCreated
		cond.Status = metav1.ConditionTrue
		meta.SetStatusCondition(&bird.Status.Conditions, *cond)
		wantUpdate = true
	}
	return wantUpdate
}

// SetupWithManager sets up the controller with the Manager.
func (r *BirdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dinov1alpha.Bird{}).
		Complete(r)
}
