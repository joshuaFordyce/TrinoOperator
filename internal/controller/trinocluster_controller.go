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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	trinov1alpha1 "github.com/joshuaFordyce/TrinoOperator/api/v1alpha1"
)

// TrinoClusterReconciler reconciles a TrinoCluster object
type TrinoClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=trino.trino.io,resources=trinoclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=trino.trino.io,resources=trinoclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=trino.trino.io,resources=trinoclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TrinoCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *TrinoClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO(user): your logic here
	
	// Fetch the TrinoCluster instance

    
	trinocluster := &trinov1alpha1.TrinoCluster{}
	err := r.Get(ctx, req.NamespacedName, trinocluster)
	
	if err != nil {
		
		
		return ctrl.Result{}, client.IgnoreNotFound(err) 
	}

			
	
		

	
	
	if trinocluster.Spec.Replicas < 1{
		logf.FromContext(ctx).Info("TrinoCluster has no replicas, skipping reconciliation", "name", trinocluster.Name, "namespace", trinocluster.Namespace)
		return ctrl.Result{}, nil
	} else {
		logf.FromContext(ctx).Info("TrinoCluster has replicas, proceeding with reconciliation", "name", trinocluster.Name, "namespace", trinocluster.Namespace)
	}
	// -- COORDINATOR DEPLOYMENT RECONCILIATION LOGIC --
	
	desiredCoordinatorDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trinocluster.Name + "-coordinator",
			Namespace: trinocluster.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func() *int32 { r := int32(1: return &r) }(), // Set the number of replicas for the coordinator deployment
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": trinocluster.Name + "-coordinator",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": trinocluster.Name + "-coordinator",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "trino-coordinator",
							Image: trinocluster.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: trinocluster.Spec.Port,
									Name:          "http",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "TRINO_CLUSTER_NAME",
									Value: trinocluster.Name + "-coordinator",
								},
								{
									Name:  "TRINO_NODE_TYPE",
									Value: "coordinator",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("1"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
						},
					},
					NodeSelector: trinocluster.Spec.NodeSelector,
					Tolerations:  trinocluster.Spec.Tolerations,
					Affinity:     trinocluster.Spec.Affinity,
				},
			},
		},
	}

    // set the OwnerReference for the coordinator deployment
	
	

	err := ctrl.SetControllerReference(trinoCluster, desiredCoordinatorDeployment, r.Scheme)
	if err != nil {
		logf.FromContext(ctx).Error(err, "Failed to set owner reference for coordinator deployment", "name", desiredCoordinatorDeployment.Name, "namespace", desireCoordinatorDeployment.Namespace)
		return ctrl.Result{}, err
	}

	foundCoordinatorDeployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: desiredCoordinatorDeployment.Name, Namespace: desiredCoordinatorDeployment.Namespace}, foundCoordinaorDeployment)

	if err != nil && errors.IsNotFound(err){
		// If the coordinator deployment is not found, create it
		logf.FromContext(ctx).Info("Coordinator deployment not found, creating", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
		err = r.Client.Create(ctx, desiredCoordinatorDeployment)
		if err != nil {
			logf.FromContext(ctx).Error(err, "Failed to create coordinator deployment", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
			return ctrl.Result{}, err
		}
		logf.FromContext(ctx).Info("Coordinator deployment created successfully", "name", desireCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
		return ctrl.Result{}, nil
		// if the coordinator deployment is not found and its a real error, log the error and return
	} else if err != nil && !errors.IsNotFound(err){
		logf.FromContext(ctx).Error(err, "Failed to get coordinator deployment", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
		return ctrl.Result{}, err
		// If the coordinator deployment is found, update it if necessary
	} else {
		logf.FromContext(ctx).Info("Coordinator deployment found, updating", "name", desireCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
		if foundCoordinatorDeployment.Spec.Replicas != desiredCoordinatorDeployment.Spec.Replicas {
			// If the replicas count has changed, update the deployment
			foundCoordinatorDeployment.Spec.Replicas = desiredCoordinatorDeployment.Spec.Replicas
			// Log the update action
			logf.FromContext(ctx).Info("Updating coordinator deployment replicas", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
			// Update the found deployment with the
			err = r.Client.Update(ctx, foundCoordinatorDeployment)
			if err != nil {
				logf.FromContext(ctx).Error(err, "Failed to update coordinator deployment", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
				return ctrl.Result{}, err
			}
			logf.FromContext(ctx).Info("Coordinator deployment updated successfully", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
		}
		// we were able to find the coordinator deployment and we didn't need to do anything
		logf.FromContext(ctx).Info("Coordinator deployment already exists, no update needed", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
		return ctrl.Result{}, nil
	}


	// -- WORKER DEPLOYMENT RECONCILIATION LOGIC --
	
	
	workerDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trinocluster.Name + "-worker",
			Namespace: trinocluster.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &trinocluster.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": trinocluster.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": trinocluster.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "trino-worker",
							Image: trinocluster.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: trinocluster.Spec.Port,
									Name:          "http",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "TRINO_CLUSTER_NAME",
									Value: trinocluster.Name,
								},
								{
									Name:  "TRINO_NODE_TYPE",
									Value: "worker",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("1"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
						},
					},
					NodeSelector: trinocluster.Spec.NodeSelector,
					Tolerations:  trinocluster.Spec.Tolerations,
					Affinity:     trinocluster.Spec.Affinity,
				},
			},
		},
	}
	
	// Create or update the worker deployment

	// Set the OwnerReference for the worker deployment
	err = ctrl.SetControllerReference(trinocluster, workerDeployment, r.Scheme)
	if err != nil {
		logf.FromContext(ctx).Error(err, "Failed to set owner reference for worker deployment", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
		return ctrl.Result{}, err
	}
	// Check if the worker deployment already exists
	foundworkerDeployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: workerDeployment.Name, Namespace: workerDeployment.Namespace}, foundworkerDeployment)

	if err != nil && errors.IsNotFound(err){
		// If the coordinator deployment is not found, create it
		logf.FromContext(ctx).Info("worker deployment not found, creating", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
		err = r.Client.Create(ctx, workerDeployment)
		if err != nil {
			logf.FromContext(ctx).Error(err, "Failed to create worker deployment", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
			return ctrl.Result{}, err
		}
		logf.FromContext(ctx).Info("worker deployment created successfully", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
		return ctrl.Result{}, nil
		// if the coordinator deployment is not found and its a real error, log the error and return
	} else if err != nil && !errors.IsNotFound(err){
		logf.FromContext(ctx).Error(err, "Failed to get worker deployment", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
		return ctrl.Result{}, err
		// If the coordinator deployment is found, update it if necessary
	} else {
		logf.FromContext(ctx).Info("worker deployment found, updating", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
		if foundworkerDeployment.Spec.Replicas != workerDeployment.Spec.Replicas {
			// If the replicas count has changed, update the deployment
			foundworkerDeployment.Spec.Replicas = workerDeployment.Spec.Replicas
			// Log the update action
			logf.FromContext(ctx).Info("Updating worker deployment replicas", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
			// Update the found deployment with the
			err = r.Client.Update(ctx, workerDeployment)
			if err != nil {
				logf.FromContext(ctx).Error(err, "Failed to update worker deployment", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
				return ctrl.Result{}, err
			}
			logf.FromContext(ctx).Info("worker deployment updated successfully", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
		}
		// we were able to find the coordinator deployment and we didn't need to do anything
		logf.FromContext(ctx).Info("worker deployment already exists, no update needed", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
		return ctrl.Result{}, nil
	}

	
	

	// ---TRIONOCLUSTER LEVEL RECONCILIATION LOGIC ---
	
	
	
    // -- Container Image update logic --

	/// -- CONFIGMAP and SECRET RECONCILILIATION LOGIC --

	// -- NETWORK POLICY, SERVICE, and INGRESS RECONCILIATION LOGIC --


}

// SetupWithManager sets up the controller with the Manager.
func (r *TrinoClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&trinov1alpha1.TrinoCluster{}).
		Named("trinocluster").
		Complete(r)
}
