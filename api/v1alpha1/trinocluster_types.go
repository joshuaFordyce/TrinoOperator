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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	networkingv1 "k8s.io/api/networking/v1"
)
type KedaAdvanced struct {
	HorizontalPodAutoscalerConfig map[string]string `json:"horizontalPodAutoscalerConfig,omitempty"`
}
type KedaFallback struct{
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
}
type AutoscalingSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	MinReplicas int32 `json:"minReplicas,omitempty"`
	MaxReplicas int32 `json:"maxReplicas,omitempty"`
	TargetCPUUtilizationPercentage int32 `json:"targetCPUUtilizationPercentage,omitempty"`
	TargetMemoryUtilizationPercentage int32 `json:"targetMemoryUtilizationPercentage,omitempty"`
	ScaleDownDelaySeconds int32 `json:"scaleDownDelaySeconds,omitempty"`
	ScaleUpDelaySeconds int32 `json:"scaleUpDelaySeconds,omitempty"`
	ScaleDownStabilizationWindowSeconds int32 `json:"scaleDownStabilizationWindowSeconds,omitempty"`
	ScaleUpStabilizationWindowSeconds int32 `json:"scaleUpStabilizationWindowSeconds,omitempty"`
	ScaleDownUnreadyReplicas int32 `json:"scaleDownUnreadyReplicas,omitempty"`
	ScaleDownUtilizationThreshold int32 `json:"scaleDownUtilizationThreshold,omitempty"`
	ScaleDownUtilizationThresholdPercentage int32 `json:"scaleDownUtilizationThresholdPercentage,omitempty"`
	ScaleDownUnreadyReplicasThreshold int32 `json:"scaleDownUnreadyReplicasThreshold,omitempty"`
	ScaleDownUnreadyReplicasPercentage int32 `json:"scaleDownUnreadyReplicasPercentage,omitempty"`
	ScaleDownUnreadyReplicasStabilizationWindowSeconds int32 `json:"scaleDownUnreadyReplicasStabilizationWindowSeconds,omitempty"`
	ScaleDownUnreadyReplicasStabilizationWindowPercentage int32 `json:"scaleDownUnreadyReplicasStabilizationWindowPercentage,omitempty"`
	ScaleDownUnreadyReplicasStabilizationWindowSecondsPercentage int32 `json:"scaleDownUnreadyReplicasStabilizationWindowSecondsPercentage,omitempty"`
}
// KedaTrigger defines the configuration for a KEDA trigger within the TrinoCluster resource.
// It specifies the type of trigger, the metric type used for scaling, and associated metadata.
// The metadata field is expected to contain trigger-specific configuration details.
type KedaTrigger struct {
	Type string `json:"type,omitempty"`
	MetricType string `json:"metricType,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`

}

// KedaSpec defines the configuration for KEDA (Kubernetes Event-driven Autoscaling) integration.
// It allows enabling KEDA-based autoscaling for the Trino cluster and configuring its behavior.
//
// Fields:
//   - Enabled: Enables or disables KEDA-based autoscaling.
//   - PollingInterval: The interval (in seconds) at which KEDA polls for metrics.
//   - CooldownPeriod: The period (in seconds) to wait after the last trigger is deactivated before scaling down.
//   - MinReplicaCount: The minimum number of replicas to maintain.
//   - MaxReplicaCount: The maximum number of replicas to allow.
//   - Fallback: Fallback configuration to use when KEDA is unavailable or fails.
//   - Advanced: Advanced KEDA configuration options.
//   - Triggers: List of KEDA triggers that define the scaling criteria.
//   - Annotations: Additional annotations to add to the KEDA ScaledObject resource.
type KedaSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	PollingInterval int32 `json:"pollingInterval,omitempty"`
	CooldownPeriod int32 `json:"cooldownPeriod,omitempty"`
	MinReplicaCount int32 `json:"minReplicaCount,omitempty"`
	MaxReplicaCount int32 `json:"maxReplicaCount,omitempty"`
	Fallback KedaFallback `json:"fallback,omitempty"`
	Advanced KedaAdvanced `json:"advanced,omitempty"`
	Triggers []KedaTrigger `json:"triggers,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

}

type ServerNodeSpec struct {
	Environment string `json:"environment,omitempty"`
	DataDir string `json:"dataDir,omitempty"`
	PluginDir string `json:"pluginDir,omitempty"`
}

type ServerLogSpec struct {
	Trino map[string]string `json:"trino,omitempty"`
}

type ServerConfigHTTPS struct {
	Enabled bool `json:"enabled,omitempty"`
	Port int32 `json:"port,omitempty"`
	Keystore map[string]string `json:"keystore,omitempty"`
}

type ServerConfigSpec struct {
	HTTP map[string]ServerConfigHTTPS `json:"http,omitempty"`
	Path string `json:"path,omitempty"`
	AuthenticationType string `json:"authenticationType,omitempty"`
	Query map[string]string `json:"query,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.


type ServerSpec struct {
	
	Node ServerNodeSpec `json:"node,omitempty"`
	Log ServerLogSpec `json:"log,omitempty"`
	Config ServerConfigSpec `json:"config,omitempty"`
	Coordinator bool `json:"coordinator,omitempty"`
	Workers bool `json:"workers,omitempty"`
	ExchangeManager map[string]string `json:"exchangeManager,omitempty"`
	CoordinatorExchangeManager map[string]string `json:"coordinatorExchangeManager,omitempty"`
	WorkerExchangeManager map[string]string `json:"workerExchangeManager,omitempty"`
	Keda KedaSpec `json:"keda,omitempty"`

}

// jmxSpec defines the configuration options for enabling and customizing JMX (Java Management Extensions)
// monitoring in the Trino cluster. It includes settings for enabling JMX, specifying the port, authentication
// credentials, and configuring the Java agent used for JMX monitoring.
//
// Fields:
//   - Enabled: Indicates whether JMX monitoring is enabled.
//   - Port: The port number on which the JMX service will be exposed.
//   - Username: The username for JMX authentication.
//   - Password: The password for JMX authentication.
//   - JavaAgentPath: The file system path to the Java agent used for JMX.
//   - JavaAgentArgs: Additional arguments to pass to the Java agent.
//   - JavaAgentEnabled: Indicates whether the Java agent for JMX is enabled.
//
// Note: There are multiple duplicate JavaAgentEnabled fields which should be consolidated into a single field.

type JmxSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	Port int32 `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	JavaAgentPath string `json:"javaAgentPath,omitempty"`
	JavaAgentArgs string `json:"javaAgentArgs,omitempty"`
	JavaAgentEnabled bool `json:"javaAgentEnabled,omitempty"`
		
}
type NetworkPolicySpec struct {
	Enabled bool `json:"enabled,omitempty"`
	PolicyTypes []networkingv1.PolicyType `json:"policyTypes,omitempty"`
	Ingress []networkingv1.NetworkPolicyIngressRule `json:"ingress,omitempty"`
	Egress []networkingv1.NetworkPolicyEgressRule `json:"egress,omitempty"`
	PodSelector metav1.LabelSelector `json:"podSelector,omitempty"`
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}
// ServiceSpec defines the configuration for the Trino service, including its type, ports, and annotations.
type ServiceSpec struct {
	Type corev1.ServiceType `json:"type,omitempty"`
	Ports []corev1.ServicePort `json:"ports,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
	ClusterIP string `json:"clusterIP,omitempty"`
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty"`
	SessionAffinity corev1.ServiceAffinity `json:"sessionAffinity,omitempty"`
}

// IngressSpec defines the configuration for ingress resources associated with the Trino cluster.
// It allows customization of ingress rules, TLS settings, backend configuration, and service exposure options.
//
// Fields:
//   - Enabled: Enables or disables ingress for the Trino cluster.
//   - ClassName: Specifies the ingress class to use.
//   - Annotations: Additional annotations to add to the ingress resource.
//   - Rules: List of ingress rules for routing external traffic.
//   - TLS: TLS configuration for securing ingress traffic.
//   - Backend: Default backend to use if no rules match.
//   - Paths: HTTP paths to expose via ingress.
//   - Host: Hostname for the ingress resource.
//   - Port: Service port to expose via ingress.
//   - TargetPort: Target port on the Trino pods.
//   - ExternalName: External DNS name for externalName services.
//   - ExternalTrafficPolicy: Controls how external traffic is routed to the service.
//   - PublishNotReadyAddresses: Whether to publish addresses of not-ready pods.
//   - ShareProcessNamespace: Whether to enable process namespace sharing for the pods.
//   - NetworkPolicy: Network policy configuration for the ingress.
type IngressSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	ClassName string `json:"className,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Rules []networkingv1.IngressRule `json:"rules,omitempty"`
	TLS []networkingv1.IngressTLS `json:"tls,omitempty"`
	Backend networkingv1.IngressBackend `json:"backend,omitempty"`
	Paths []networkingv1.HTTPIngressPath `json:"paths,omitempty"`
	Host string `json:"host,omitempty"`
	Port int32 `json:"port,omitempty"`
	TargetPort int32 `json:"targetPort,omitempty"`
	ExternalName string `json:"externalName,omitempty"`
	ExternalTrafficPolicy corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	PublishNotReadyAddresses bool `json:"publishNotReadyAddresses,omitempty"`
	ShareProcessNamespace bool `json:"shareProcessNamespace,omitempty"`
	NetworkPolicy NetworkPolicySpec `json:"networkPolicy,omitempty"`
}

type ServiceMonitorSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	Selector metav1.LabelSelector `json:"selector,omitempty"`
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	Interval string `json:"interval,omitempty"`
	Timeout string `json:"timeout,omitempty"`
	Port string `json:"port,omitempty"`
	
	Labels map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	TargetLabels []string `json:"targetLabels,omitempty"`
	MetricLabels []string `json:"metricLabels,omitempty"`
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`
	SelectorAnnotations map[string]string `json:"selectorAnnotations,omitempty"`
	SelectorMatchLabels map[string]string `json:"selectorMatchLabels,omitempty"`
	SelectorMatchExpressions []metav1.LabelSelectorRequirement `json:"selectorMatchExpressions,omitempty"`
}
type CoordinatorSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	Affinity *corev1.Affinity `json:"affinity,omitempty"`	
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	PodDisruptionBudget map[string]string `json:"podDisruptionBudget,omitempty"`
	Autoscaling AutoscalingSpec `json:"autoscaling,omitempty"`
	Advanced KedaAdvanced `json:"advanced,omitempty"`
	Fallback KedaFallback `json:"fallback,omitempty"`
	Triggers []KedaTrigger `json:"triggers,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Env []corev1.EnvVar `json:"env,omitempty"`
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
	SidecarContainers []corev1.Container `json:"sidecarContainers,omitempty"`
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
	ContainerSecurityContext *corev1.SecurityContext `json:"containerSecurityContext,omitempty"`
	ShareProcessNamespace map[string]bool `json:"shareProcessNamespace,omitempty"`
	Service ServiceSpec `json:"service,omitempty"`
	Auth string`json:"auth,omitempty"`
	ServiceAccount string `json:"serviceAccount,omitempty"`
	ConfigMounts []corev1.VolumeMount `json:"configMounts,omitempty"`
	SecretMounts []corev1.VolumeMount `json:"secretMounts,omitempty"`
}
type WorkerSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	PodDisruptionBudget map[string]string `json:"podDisruptionBudget,omitempty"`
	Autoscaling AutoscalingSpec `json:"autoscaling,omitempty"`
	Advanced KedaAdvanced `json:"advanced,omitempty"`
	Fallback KedaFallback `json:"fallback,omitempty"`
	Triggers []KedaTrigger `json:"triggers,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Env []corev1.EnvVar `json:"env,omitempty"`
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
	SidecarContainers []corev1.Container `json:"sidecarContainers,omitempty"`
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
	ContainerSecurityContext *corev1.SecurityContext `json:"containerSecurityContext,omitempty"`
	ShareProcessNamespace map[string]bool `json:"shareProcessNamespace,omitempty"`
	Service ServiceSpec `json:"service,omitempty"`
	Auth string`json:"auth,omitempty"`
	ServiceAccount string `json:"serviceAccount,omitempty"`
	ConfigMounts []corev1.VolumeMount `json:"configMounts,omitempty"`
	SecretMounts []corev1.VolumeMount `json:"secretMounts,omitempty"`
}

// ImageSpec defines the configuration for the container image used in the Trino cluster.

type ImageSpec struct {
	Repository string `json:"repository,omitempty"`
	Tag string `json:"tag,omitempty"`
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
	UseRepositoryasSoleImageReference bool `json:"useRepositoryasSoleImageReference,omitempty"`	
}

// TrinoClusterSpec defines the desired state of TrinoCluster.
type TrinoClusterSpec struct {

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of TrinoCluster. Edit trinocluster_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
	NameOverride        string `json:"nameOverride,omitempty"`
	CoordinatorNameOverride string `json:"coordinatorNameOverride,omitempty"`
	WorkerNameOverride  string `json:"workerNameOverride,omitempty"`

	Image ImageSpec `json:"image,omitempty"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	Server ServerSpec `json:"server,omitempty"`
	
	AccessControl string `json:"accessControl,omitempty"`
	ResourceGroups string `json:"resourceGroups,omitempty"`

	AdditionalNodeProperties []string `json:"additionalNodeProperties,omitempty"`
	AdditionalConfigProperties []string `json:"additionalConfigProperties,omitempty"`
	AdditionalLogProperties []string `json:"additionalLogProperties,omitempty"`
	AdditionalExchangeManagerProperties []string `json:"additionalExchangeManagerProperties,omitempty"`
	
	SessionProperties string`json:"sessionProperties,omitempty"`
	EventListenerProperties []string `json:"eventListenerProperties,omitempty"`

	Catalogs map[string]string `json:"catalogs,omitempty"`

	Env []corev1.EnvVar `json:"env,omitempty"`
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`

	InitContainers []corev1.Container `json:"initContainers,omitempty"`
	SidecarContainers []corev1.Container `json:"sidecarContainers,omitempty"`

	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
	ContainerSecurityContext *corev1.SecurityContext `json:"containerSecurityContext,omitempty"`
	ShareProcessNamespace map[string]bool `json:"shareProcessNamespace,omitempty"`

	Service ServiceSpec `json:"service,omitempty"`
	Auth string `json:"auth,omitempty"`
	ServiceAccount string`json:"serviceAccount,omitempty"`

	ConfigMounts []corev1.VolumeMount `json:"configMounts,omitempty"`
	SecretMounts []corev1.VolumeMount `json:"secretMounts,omitempty"`

	Coordinator CoordinatorSpec `json:"coordinator,omitempty"`
	Worker WorkerSpec `json:"worker,omitempty"`

	Kafka map[string]string `json:"kafka,omitempty"`
	Jmx JmxSpec `json:"jmx,omitempty"`
	ServiceMonitor ServiceMonitorSpec `json:"serviceMonitor,omitempty"`
	
	CommonLabels map[string]string `json:"commonLabels,omitempty"`

	Ingress IngressSpec `json:"ingress,omitempty"`
	NetworkPolicy NetworkPolicySpec `json:"networkPolicy,omitempty"`
	}

// TrinoClusterStatus defines the observed state of TrinoCluster.
type TrinoClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	CoordinatorStatus appsv1.StatefulSetStatus `json:"coordinatorStatus,omitempty"`
	WorkerStatus      appsv1.StatefulSetStatus `json:"workerStatus,omitempty"`
	WorkerCount      int32                   `json:"workerCount,omitempty"`
	ReadyWorkerCount int32 `json:"readyWorkerCount,omitempty"`
	ClusterState string `json:"clusterState,omitempty"`
	TrinoUIURL string `json:"trinoUIURL,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	TrinoVersion string `json:"trinoVersion,omitempty"`
	ConfigHash string `json:"configHash,omitempty"`
	// Add more fields as needed to track the status of the Trino cluster
}


// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TrinoCluster is the Schema for the trinoclusters API.
type TrinoCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrinoClusterSpec   `json:"spec,omitempty"`
	Status TrinoClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TrinoClusterList contains a list of TrinoCluster.
type TrinoClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrinoCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrinoCluster{}, &TrinoClusterList{})
}
