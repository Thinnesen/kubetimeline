/*
Copyright 2025.

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
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	timelinev1alpha1 "github.com/Thinnesen/kubetimeline/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// KubeTimelineReconciler reconciles a KubeTimeline object
type KubeTimelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=timeline.thinnesen.com,resources=kubetimelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=timeline.thinnesen.com,resources=kubetimelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=timeline.thinnesen.com,resources=kubetimelines/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KubeTimeline object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile

// Helper to get the best timestamp for an event
func getEventTimestamp(e corev1.Event) time.Time {
	if !e.LastTimestamp.IsZero() {
		return e.LastTimestamp.Time
	}
	if !e.EventTime.IsZero() {
		return e.EventTime.Time
	}
	return e.CreationTimestamp.Time
}

func (r *KubeTimelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the KubeTimeline instance
	var timeline timelinev1alpha1.KubeTimeline
	if err := r.Get(ctx, req.NamespacedName, &timeline); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var eventList corev1.EventList
	var listOpts []client.ListOption

	if timeline.Spec.ClusterWide {
		listOpts = []client.ListOption{}
	} else {
		listOpts = []client.ListOption{client.InNamespace(req.Namespace)}
	}

	if err := r.List(ctx, &eventList, listOpts...); err != nil {
		log.Error(err, "unable to list events for KubeTimeline")
		return ctrl.Result{}, err
	}

	// Sort events by timestamp descending (latest first)
	sort.Slice(eventList.Items, func(i, j int) bool {
		return getEventTimestamp(eventList.Items[i]).After(getEventTimestamp(eventList.Items[j]))
	})

	var messages []string
	for i := 0; i < len(eventList.Items) && len(messages) < 25; i++ {
		e := eventList.Items[i]
		timestamp := getEventTimestamp(e).String()
		messages = append(messages, timestamp+" "+e.InvolvedObject.Kind+"/"+e.InvolvedObject.Name+": "+e.Message+"🐒")
	}

	if !equalStringSlices(timeline.Status.Events, messages) {
		timeline.Status.Events = messages
		if err := r.Status().Update(ctx, &timeline); err != nil {
			log.Error(err, "unable to update KubeTimeline status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeTimelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&timelinev1alpha1.KubeTimeline{}).
		Named("kubetimeline").
		Complete(r)
}
