/*
Copyright 2025.
...
*/

package controller

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	trinov1alpha1 "github.com/joshuaFordyce/TrinoOperator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
    
    serviceName := fmt.Sprintf("%s-coordinator-service", trinoCluster.Name)

    servicePort := trinoCluster.Spec.Coordinator.Service.Ports[0].Port

    return &corev1.Service {
        ObjectMeta: metav1.ObjectMeta{
            Name: serviceName,
            Namespace: trinoCluster.Namespace,
            Labels: map[string]string{"app": trinoCluster.Name},
        },
        Spec: corev1.ServiceSpec{

            Selector: map[string]string{
                "app": trinoCluster.Name,
                "app.kubernetes.io/component": "coordinator",
            },
            Ports: []corev1.ServicePort {
                {
                    Protocol: corev1.ProtocolTCP,
                    Port:     servicePort,
                    TargetPort: intstr.FromInt(int(servicePort)),
                    Name: "http",
                },
            },
           
            Type: corev1.ServiceTypeClusterIP,
        },
    }
}

func (r *TrinoClusterReconciler) reconcileCoordinatorService( ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {
    
    log := logf.FromContext((ctx))

    desiredCoordinatorService := r.getCoordinatorService(trinoCluster)

    if err := ctrl.SetControllerReference(trinoCluster, desiredCoordinatorService, r.Scheme); err != nil{
        return fmt.Errorf("failed to set owner reference: %w", err)
    }

    foundCoordinatorService := &corev1.Service{}
    err := r.Client.Get(ctx, types.NamespacedName{Name: desiredCoordinatorService.Name, Namespace: desiredCoordinatorService.Namespace}, foundCoordinatorService)

    if err != nil && errors.IsNotFound(err) {
        //Create: The service was not found, so we create it
        log.Info("Creating the Service for the Coordinator Nodes", "name", desiredCoordinatorService.Name)
        return r.Client.Create(ctx, desiredCoordinatorService)
    } else if err != nil {
        // ERROR: a real error happend so return it
        return fmt.Errorf("failed to get coordinator service: %w", err)
    }

    //UPDATE: The service was found check if update needed
    if !reflect.DeepEqual(foundCoordinatorService.Spec, desiredCoordinatorService.Spec) {
        log.Info("Updating coordinator service spec", "name", foundCoordinatorService.Name)
        foundCoordinatorService.Spec = desiredCoordinatorService.Spec
        if err = r.Client.Update(ctx, foundCoordinatorService); err != nil {
            log.Error(err, "Failed to update the service for the Coordinator", "name",foundCoordinatorService.Name)
            return fmt.Errorf("failed to update Coordinator service:  %w", err)
        }
    }
    log.Info("Coordinator service reconciled successfully", "name", foundCoordinatorService.Name)
    return nil

}

func (r *TrinoClusterReconciler) getWorkerService(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.Service {
    serviceName := fmt.Sprintf("%s-worker-service", trinoCluster.Name)

    servicePort := trinoCluster.Spec.Coordinator.Service.Ports[0].Port

    return &corev1.Service {
        ObjectMeta: metav1.ObjectMeta{
            Name: serviceName,
            Namespace: trinoCluster.Namespace,
            Labels: map[string]string{"app": trinoCluster.Name},
        },
        Spec: corev1.ServiceSpec{

            Selector: map[string]string{
                "app": trinoCluster.Name,
                "app.kubernetes.io/component": "worker",
            },
            Ports: []corev1.ServicePort {
                {
                    Protocol: corev1.ProtocolTCP,
                    Port:     servicePort,
                    TargetPort: intstr.FromInt(int(servicePort)),
                    Name: "http",
                },
            },
           
            Type: corev1.ServiceTypeClusterIP,
        },
    }    

}
func (r *TrinoClusterReconciler) reconcileWorkerService( ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {

    log := logf.FromContext(ctx)
    desiredWorkerService := r.getWorkerService(trinoCluster)

    if err := ctrl.SetControllerReference(trinoCluster, desiredWorkerService, r.Scheme); err != nil{
        return err
    }

    foundWorkerService := &corev1.Service{}
    err := r.Client.Get(ctx, types.NamespacedName{Name: desiredWorkerService.Name, Namespace: desiredWorkerService.Namespace}, foundWorkerService)

   if err != nil && errors.IsNotFound(err) {
        //Create: The service was not found, so we create it
        log.Info("Creating the Service for the Worker Nodes", "name", desiredWorkerService.Name)
        return r.Client.Create(ctx, desiredWorkerService)
    } else if err != nil {
        // ERROR: a real error happend so return it
        return fmt.Errorf("failed to get coordinator service: %w", err)
    }

    //UPDATE: The service was found check if update needed
    if !reflect.DeepEqual(foundWorkerService.Spec, desiredWorkerService.Spec) {
        log.Info("Updating coordinator service spec", "name", foundWorkerService.Name)
        foundWorkerService.Spec = desiredWorkerService.Spec
        if err = r.Client.Update(ctx, foundWorkerService); err != nil {
            log.Error(err, "Failed to update the service for the Worker", "name",foundWorkerService.Name)
            return fmt.Errorf("failed to update Worker service:  %w", err)
        }
    }
    log.Info("Worker service reconciled successfully", "name", foundWorkerService.Name)
    return nil
}

func (r *TrinoClusterReconciler) getWorkerCM(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.ConfigMap {
    configMapName := fmt.Sprintf("%s-workerCM", trinoCluster.Name)

    var nodeProperties []string
    nodeProperties = append(nodeProperties,"node.environment=production")
    nodeProperties = append(nodeProperties,"node.data-dir=/data/trino")

    nodePropertiesString := strings.Join(nodeProperties, "\n")

    

    var jvmConfig []string

    jvmConfig = append(jvmConfig, "-Xmx8g")
    jvmConfig = append(jvmConfig, "-XX:+UseG1GC")
    jvmConfig = append(jvmConfig, "-XX:G1HeapRegionSize=32M")

    jvmConfigString := strings.Join(jvmConfig, "\n")

    var configProperties []string
    configProperties = append(configProperties, "coordinator=false")
    configProperties = append(configProperties, "http-server.http.port=8080")
    for key, value := range trinoCluster.Spec.Worker.Config.Properties {
        configProperties = append(configProperties, fmt.Sprintf("%s=%s", key, value))
    }
    

    configPropertiesString := strings.Join(configProperties, "\n")

    return &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: configMapName,
            Namespace: trinoCluster.Namespace,
            Labels: map[string]string{"app": trinoCluster.Name},
        },
        Data: map[string]string{
            "node.properties":nodePropertiesString,
            "jvm.config": jvmConfigString,
            "config.properties": configPropertiesString,
        },
    }

}

func (r *TrinoClusterReconciler) reconcileWorkerCM(ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {
    
    log := logf.FromContext(ctx)

    desiredWorkerCM := r.getWorkerCM(trinoCluster)

    if err := ctrl.SetControllerReference(trinoCluster, desiredWorkerCM, r.Scheme); err != nil {
        return err
    }

    foundConfigMap := &corev1.ConfigMap{}

    err := r.Client.Get(ctx, types.NamespacedName{Name: desiredWorkerCM.Name, Namespace: desiredWorkerCM.Namespace}, foundConfigMap)

    if err != nil && errors.IsNotFound(err) {
        log.Info("Creating worker ConfigMap", "name", desiredWorkerCM.Name)
        return err
    }

    if !reflect.DeepEqual(foundConfigMap.Data, desiredWorkerCM.Data) {
        foundConfigMap.Data = desiredWorkerCM.Data
        if err = r.Client.Update(ctx, foundConfigMap); err != nil {
            log.Error(err, "Failed to update ConfigMap", "name", foundConfigMap.Name)
            return fmt.Errorf("this is the err from not being able to create a configmap %w",err)
        }
    }
    return nil

}
func (r *TrinoClusterReconciler) getCoordinatorCM(trinoCluster *trinov1alpha1.TrinoCluster) *corev1.ConfigMap {
   configMapName := fmt.Sprintf("%s-coordinatorCM", trinoCluster.Name)

    var nodeProperties []string
    nodeProperties = append(nodeProperties,"node.environment=production")
    nodeProperties = append(nodeProperties,"node.data-dir=/data/trino")

    nodePropertiesString := strings.Join(nodeProperties, "\n")

    

    var jvmConfig []string

    jvmConfig = append(jvmConfig, "-Xmx8g")
    jvmConfig = append(jvmConfig, "-XX:+UseG1GC")
    jvmConfig = append(jvmConfig, "-XX:G1HeapRegionSize=32M")

    jvmConfigString := strings.Join(jvmConfig, "\n")

    var configProperties []string
    configProperties = append(configProperties, "coordinator=true")
    configProperties = append(configProperties, "http-server.http.port=8080")
    for key, value := range trinoCluster.Spec.Worker.Config.Properties {
        configProperties = append(configProperties, fmt.Sprintf("%s=%s", key, value))
    }
    

    configPropertiesString := strings.Join(configProperties, "\n")

    return &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: configMapName,
            Namespace: trinoCluster.Namespace,
            Labels: map[string]string{"app": trinoCluster.Name},
        },
        Data: map[string]string{
            "node.properties":nodePropertiesString,
            "jvm.config": jvmConfigString,
            "config.properties": configPropertiesString,
        },
    }

}

func (r *TrinoClusterReconciler) reconcileCoordinatorCM(ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {

    log := logf.FromContext(ctx)

    desiredCoordinatorCM := r.getCoordinatorCM(trinoCluster)

    if err := ctrl.SetControllerReference(trinoCluster, desiredCoordinatorCM, r.Scheme); err != nil {
        return err
    }

    foundConfigMap := &corev1.ConfigMap{}

    err := r.Client.Get(ctx, types.NamespacedName{Name: desiredCoordinatorCM.Name, Namespace: desiredCoordinatorCM.Namespace}, foundConfigMap)

    if err != nil && errors.IsNotFound(err) {
        log.Info("Creating Coordinator ConfigMap", "name", desiredCoordinatorCM.Name)
        return err
    }

    if !reflect.DeepEqual(foundConfigMap.Data, desiredCoordinatorCM.Data) {
        foundConfigMap.Data = desiredCoordinatorCM.Data
        if err = r.Client.Update(ctx, foundConfigMap); err != nil {
            log.Error(err, "Failed to update ConfigMap", "name", foundConfigMap.Name)
            return fmt.Errorf("failed to updateConfigMap %w", err)
        }
    }
    return nil


}


func (r *TrinoClusterReconciler) getCatalogCM(trinoCluster *trinov1alpha1.TrinoCluster, catalogName string, catalogSpec trinov1alpha1.CatalogSpec) *corev1.ConfigMap {
    var props []string
    for key, value := range catalogSpec.Properties {
        props = append(props, fmt.Sprintf("%s=%s", key, value))
    }
    propertiesString := strings.Join(props, "\n")

    return &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:  fmt.Sprintf("%s-catalog-%s", trinoCluster.Name, catalogName),
            Namespace: trinoCluster.Namespace,
            Labels: map[string]string{"app": trinoCluster.Name},
        },
        Data: map[string]string{
            fmt.Sprintf("%s.properties", catalogName): propertiesString,
        },

    }

}

func (r *TrinoClusterReconciler) reconcileCatalogCM(ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {

    log := logf.FromContext(ctx)

    for catalogName, catalogSpec := range trinoCluster.Spec.Catalogs {
        var err error

        desiredConfigMap := r.getCatalogCM(trinoCluster, catalogName, catalogSpec)

        if err := ctrl.SetControllerReference(trinoCluster, desiredConfigMap, r.Scheme); err != nil {
            return err
        }

        foundConfigMap := &corev1.ConfigMap{}
        err = r.Client.Get(ctx, types.NamespacedName{Name: desiredConfigMap.Name, Namespace: desiredConfigMap.Namespace}, foundConfigMap)

        if err != nil && errors.IsNotFound(err) {

            log.Info("Creating catalog ConfigMap", "name", desiredConfigMap.Name)
            if err = r.Client.Create(ctx, desiredConfigMap); err != nil {
                
                return fmt.Errorf("failed to create catalog ConfigMap for %s: %w", catalogName, err)

            }
            
        } else if err != nil && !errors.IsNotFound(err) {
            return fmt.Errorf("failed to get ConfigMap for %s: %w", catalogName,err)
        }else {

        if !reflect.DeepEqual(foundConfigMap.Data, desiredConfigMap.Data) {
            foundConfigMap.Data = desiredConfigMap.Data
            if err := r.Client.Update(ctx, foundConfigMap); err != nil {
                log.Error(err, "Failed to update ConfigMap", "name", foundConfigMap.Name)
                return fmt.Errorf("failed to update ConfigMap for %s: %w", catalogName, err)
            
                }
        
            }
        }
    
    }
    return nil
}


func (r *TrinoClusterReconciler) getIngress(trinoCluster *trinov1alpha1.TrinoCluster) *networkingv1.Ingress {
    ingressName := fmt.Sprintf("%s-ingress", trinoCluster.Name)
    
    pathType := networkingv1.PathTypePrefix // This is a required field.
    
    return &networkingv1.Ingress{
        ObjectMeta: metav1.ObjectMeta{
            Name: ingressName,
            Namespace: trinoCluster.Namespace,
            Labels: map[string]string{
                "app": trinoCluster.Name,
            },
            Annotations: map[string]string{
                // Add any Ingress-specific annotations here (e.g., for NGINX or AWS).
                "nginx.ingress.kubernetes.io/rewrite-target": "/",
            },
        },
        Spec: networkingv1.IngressSpec{
            Rules: []networkingv1.IngressRule{
                {
                    Host: trinoCluster.Spec.Ingress.Host,
                    IngressRuleValue: networkingv1.IngressRuleValue{
                        HTTP: &networkingv1.HTTPIngressRuleValue{
                            Paths: []networkingv1.HTTPIngressPath{
                                {
                                    Path: "/",
                                    PathType: &pathType,
                                    Backend: networkingv1.IngressBackend{
                                        Service: &networkingv1.IngressServiceBackend{
                                            Name: fmt.Sprintf("%s-coordinator-service", trinoCluster.Name),
                                            Port: networkingv1.ServiceBackendPort{
                                                Number: trinoCluster.Spec.Ingress.Port,
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }
}


func (r *TrinoClusterReconciler) reconcileIngress(ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {

    log := logf.FromContext(ctx)

    desiredIngress := r.getIngress(trinoCluster)
    if err := ctrl.SetControllerReference(trinoCluster, desiredIngress, r.Scheme); err != nil {
        return err
    }
    
    foundIngress := &networkingv1.Ingress{}
    err := r.Client.Get(ctx, types.NamespacedName{Name: desiredIngress.Name, Namespace: desiredIngress.Namespace}, foundIngress)
    if err != nil && errors.IsNotFound(err) {
        log.Info("Ingress is not found", "name", desiredIngress.Name, "namespace", desiredIngress.Namespace)
        return r.Client.Create(ctx, desiredIngress)
    } else if err != nil && !errors.IsNotFound(err) {
        return fmt.Errorf("failed to get Ingress for %s: %w", desiredIngress.Name, err)
    } else {
        if !reflect.DeepEqual(foundIngress.Spec, desiredIngress.Spec) {
            foundIngress.Spec = desiredIngress.Spec

            logf.FromContext(ctx).Info("Updating Ingres Spec", "name", desiredIngress.Name, "namespace", desiredIngress.Namespace)

            err = r.Client.Update(ctx, foundIngress)
            if err != nil {
                logf.FromContext(ctx).Error(err, "Failed to update Ingress spec", "name", desiredIngress.Name, "namespace", desiredIngress.Namespace)
                return err
            }
        }
        return nil
    }
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
                logf.FromContext(ctx).Error(err, "Failed to update worker spec", "name", foundworkerDeployment.Name, "namespace", foundworkerDeployment.Namespace)
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

func (r *TrinoClusterReconciler) cleanup(ctx context.Context, trinoCluster *trinov1alpha1.TrinoCluster) error {
    // Delete the coordinator deployment

    if err := r.Client.Delete(ctx, &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-coordinator", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

    if err := r.Client.Delete(ctx, &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-worker", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

    if err := r.Client.Delete(ctx, &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-coordinator-service", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

    if err := r.Client.Delete(ctx, &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-worker-service", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

    if err := r.Client.Delete(ctx, &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-workerCM", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

    if err := r.Client.Delete(ctx, &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-coordinatorCM", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

    if err := r.Client.Delete(ctx, &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-catalog", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

    if err := r.Client.Delete(ctx, &networkingv1.Ingress{
        ObjectMeta: metav1.ObjectMeta{
            Name: fmt.Sprintf("%s-ingress", trinoCluster.Name),
            Namespace: trinoCluster.Namespace,
        },
    }); client.IgnoreNotFound(err) != nil {
        return err
    }

return nil
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

   

    

    /// -- CONFIGMAP and SECRET RECONCILILIATION LOGIC --
    if err = r.reconcileWorkerCM(ctx, trinocluster); err != nil {
        return ctrl.Result{}, err
    }

    if err = r.reconcileCoordinatorCM(ctx, trinocluster); err != nil {
        return ctrl.Result{}, err
    }
    
    if err = r.reconcileCatalogCM(ctx, trinocluster); err != nil {
        return ctrl.Result{}, err
    }


    // -- NETWORK POLICY, SERVICE, and INGRESS RECONCILIATION LOGIC --
     //resolve services and ingresses

    if err = r.reconcileCoordinatorService(ctx, trinocluster); err != nil {
        return ctrl.Result{}, err
    }

    if err = r.reconcileWorkerService(ctx, trinocluster); err != nil {
        return ctrl.Result{}, err
    }

    if err = r.reconcileIngress(ctx,trinocluster); err != nil {
        return ctrl.Result{}, err
    }

    //Add the Finalizer logic

    const trinoClusterFinalizer = "trino.trino.io/finalizer"

    if !controllerutil.ContainsFinalizer(trinocluster, trinoClusterFinalizer) {
        //add the finalizer to the object
        controllerutil.AddFinalizer(trinocluster, trinoClusterFinalizer)
        // Update the object to save the finalizer
        if err := r.Update(ctx, trinocluster); err != nil {
            return ctrl.Result{}, err

        }



    }

    isTrinoClusterMarkedForDeleteion := trinocluster.GetDeletionTimestamp() != nil 
    if isTrinoClusterMarkedForDeleteion {
        if controllerutil.ContainsFinalizer(trinocluster, trinoClusterFinalizer) {
            // Run the cleanup logic here
            // Delete all the owned resources (deployments, services, Configmaps)
            if err := r.cleanup(ctx, trinocluster); err != nil {
                return ctrl.Result{}, err
            }

            //remove the finalizer from the object
            controllerutil.RemoveFinalizer(trinocluster, trinoClusterFinalizer)
            if err := r.Update(ctx, trinocluster); err != nil {
                return ctrl.Result{}, err
            }
        }
    }
return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TrinoClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&trinov1alpha1.TrinoCluster{}).
        Named("trinocluster").
        Complete(r)
}