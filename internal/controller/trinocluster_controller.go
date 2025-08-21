/*
Copyright 2025.
...
*/

package controller

import (
    "context"
    "reflect"
    "fmt" 
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    logf "sigs.k8s.io/controller-runtime/pkg/log"

    trinov1alpha1 "github.com/joshuaFordyce/TrinoOperator/api/v1alpha1"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/apimachinery/pkg/api/errors"
    networkingv1 "k8s.io/api/networking/v1"
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
// Helpers for the ConfigMaps
func (r *TrinoClusterReconciler) getCoordinatorService(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.Service{

}

func (r *TrinoClusterReconciler) getWorkerService(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.Service {
    

}

func (r *TrinoClusterReconciler) getWorkerCM(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.ConfigMap {

}
func (r *TrinoClusterReconciler) getCoordinatorCM(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.ConfigMap {

}

func (r *TrinoClusterReconciler) getCatalogCM(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.ConfigMap {

}

func (r *TrinoClusterReconciler) getIngress(trinoCluster *trinov1alpha1.TrinoCluster) *networkingv1.Ingress {

}

func (r *TrinoClusterReconciler) getCatalogsCM(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.ConfigMap {

}



//Helpers for the Services
// Helpers for the Ingress
// Update Status: Update the TrinoCluster CR's status field with the number of ready pods, current state, 
// Add Finalizers whihc ensure proper cleanup
func (r *TrinoClusterReconciler ) getDesiredCoordinatorDeployment(trinoCluster *trinov1alpha1.TrinoCluster) *appsv1.Deployment {
	return &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      trinoCluster.Name + "-coordinator",
            Namespace: trinoCluster.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: int32Ptr(1),
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app": trinoCluster.Name + "-coordinator",
                },
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app": trinoCluster.Name + "-coordinator",
                    },
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "trino-coordinator",
                            Image: fmt.Sprintf("%s:%s", trinoCluster.Spec.Image.Repository, trinoCluster.Spec.Image.Tag),
                            Ports: []corev1.ContainerPort{
                                {
                                    ContainerPort: trinoCluster.Spec.Coordinator.Service.Ports[0].Port,
                                    Name:          "http",
                                },
                            },
                            Env: []corev1.EnvVar{
                                {
                                    Name:  "TRINO_CLUSTER_NAME",
                                    Value: trinoCluster.Name + "-coordinator",
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
                    NodeSelector: trinoCluster.Spec.Coordinator.NodeSelector,
                    Tolerations:  trinoCluster.Spec.Coordinator.Tolerations,
                    Affinity:     trinoCluster.Spec.Coordinator.Affinity,
                },
            },
        },
    }

}

func (r *TrinoClusterReconciler) getDesiredWorkerDeployment(trinoCluster *trinov1alpha1.TrinoCluster) *appsv1.Deployment {
	return &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      trinoCluster.Name + "-worker",
            Namespace: trinoCluster.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &trinoCluster.Spec.Worker.Replicas,
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app": trinoCluster.Name,
                },
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app": trinoCluster.Name,
                    },
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "trino-worker",
                            Image: fmt.Sprintf("%s:%s", trinoCluster.Spec.Image.Repository, trinoCluster.Spec.Image.Tag),
                            Ports: []corev1.ContainerPort{
                                {
                                    ContainerPort: trinoCluster.Spec.Worker.Service.Ports[0].Port,
                                    Name:          "http",
                                },
                            },
                            Env: []corev1.EnvVar{
                                {
                                    Name:  "TRINO_CLUSTER_NAME",
                                    Value: trinoCluster.Name,
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
                    NodeSelector: trinoCluster.Spec.Worker.NodeSelector,
                    Tolerations:  trinoCluster.Spec.Worker.Tolerations,
                    Affinity:     trinoCluster.Spec.Worker.Affinity,
                },
            },
        },
    }

}

func (r *TrinoClusterReconciler ) reconcileWorkerDeployment(ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {
    workerDeployment := r.getDesiredWorkerDeployment(trinoCluster)
    

    // Set the OwnerReference for the worker deployment
    err := ctrl.SetControllerReference(trinoCluster, workerDeployment, r.Scheme)
    if err != nil {
        logf.FromContext(ctx).Error(err, "Failed to set owner reference for worker deployment", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
        return err
    }

	foundworkerDeployment := &appsv1.Deployment{}
    err = r.Client.Get(ctx, types.NamespacedName{Name: workerDeployment.Name, Namespace: workerDeployment.Namespace}, foundworkerDeployment)
    if err != nil && errors.IsNotFound(err){
        // If the worker deployment is not found, create it
        logf.FromContext(ctx).Info("worker deployment not found, creating", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
        return r.Client.Create(ctx, workerDeployment)
        
       // if the worker deployment is not found and its a real error, log the error and return
    } else if err != nil && !errors.IsNotFound(err){
        logf.FromContext(ctx).Error(err, "Failed to get worker deployment", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
        return err
    } else {
        logf.FromContext(ctx).Info("worker deployment found, updating", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
        if !reflect.DeepEqual(foundworkerDeployment.Spec, workerDeployment.Spec) {
            // If the spec has changed, update the deployment
            foundworkerDeployment.Spec = workerDeployment.Spec
            // Log the update action
            logf.FromContext(ctx).Info("Updating worker deployment spec", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)

            // Update the found deployment with the
            err = r.Client.Update(ctx, foundworkerDeployment)
            if err != nil {
                logf.FromContext(ctx).Error(err, "Failed to update worker spec", "name", workerDeployment.Name, "namespace", workerDeployment.Namespace)
                return err
            }
        }
    return nil
    }
}


func (r *TrinoClusterReconciler ) reconcileCoordinatorDeployment(ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {
	desiredCoordinatorDeployment := r.getDesiredCoordinatorDeployment(trinoCluster)

    // set the OwnerReference for the coordinator deployment
    err := ctrl.SetControllerReference(trinoCluster, desiredCoordinatorDeployment, r.Scheme)
    if err != nil {
        logf.FromContext(ctx).Error(err, "Failed to set owner reference for coordinator deployment", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
        return err
    }
	foundCoordinatorDeployment := &appsv1.Deployment{}
    err = r.Client.Get(ctx, types.NamespacedName{Name: desiredCoordinatorDeployment.Name, Namespace: desiredCoordinatorDeployment.Namespace}, foundCoordinatorDeployment)

    if err != nil && errors.IsNotFound(err){
        // If the coordinator deployment is not found, create it
        logf.FromContext(ctx).Info("Coordinator deployment not found, creating", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
        return r.Client.Create(ctx, desiredCoordinatorDeployment)
        
        // No return here, so the reconciliation can continue
    } else if err != nil {
        logf.FromContext(ctx).Error(err, "Failed to get coordinator deployment", "name", desiredCoordinatorDeployment.Name, "namespace", desiredCoordinatorDeployment.Namespace)
        return err
    } else {
        logf.FromContext(ctx).Info("Coordinator deployment found, updating", "name", foundCoordinatorDeployment.Name, "namespace", foundCoordinatorDeployment.Namespace)
        if !reflect.DeepEqual(foundCoordinatorDeployment.Spec, desiredCoordinatorDeployment.Spec) {
            // If the spec has changed, update the deployment
            foundCoordinatorDeployment.Spec = desiredCoordinatorDeployment.Spec
            logf.FromContext(ctx).Info("Updating coordinator deployment spec", "name", foundCoordinatorDeployment.Name, "namespace", foundCoordinatorDeployment.Namespace)
            err = r.Client.Update(ctx, foundCoordinatorDeployment)
            if err != nil {
                logf.FromContext(ctx).Error(err, "Failed to update coordinator deployment", "name", foundCoordinatorDeployment.Name, "namespace", foundCoordinatorDeployment.Namespace)
                return err
            }
           
        // No return here, so the reconciliation can continue
        
        
    }
    return nil

    }
}



func int32Ptr(i int32) *int32 {
    return &i
}
func (r *TrinoClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    _ = logf.FromContext(ctx)

    trinocluster := &trinov1alpha1.TrinoCluster{}
    err := r.Get(ctx, req.NamespacedName, trinocluster)
    
    if err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err) 
    }

    
    // -- COORDINATOR DEPLOYMENT RECONCILIATION LOGIC --
    
    
   

    // Check Coordinator Deployment and update if necessary
    if err := r.reconcileCoordinatorDeployment(ctx, trinocluster); err != nil {
       
        return ctrl.Result{}, fmt.Errorf("failed to re")
    }

   


    // -- WORKER DEPLOYMENT RECONCILIATION LOGIC --
    
   
    
    // Create or update the worker deployment
    
    // Check Worker Deployment and update if Necessary

    if err = r.reconcileWorkerDeployment(ctx, trinocluster); err != nil {
       
        return ctrl.Result{}, err
    }
    
    
    

    // ---TRIONOCLUSTER LEVEL RECONCILIATION LOGIC ---
    
    
    
    // -- Container Image update logic --


    /// -- CONFIGMAP and SECRET RECONCILILIATION LOGIC --

    // -- NETWORK POLICY, SERVICE, and INGRESS RECONCILIATION LOGIC --
 
return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TrinoClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&trinov1alpha1.TrinoCluster{}).
        Named("trinocluster").
        Complete(r)
}