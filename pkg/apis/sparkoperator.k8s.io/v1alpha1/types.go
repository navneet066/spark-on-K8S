/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SparkApplicationType describes the type of a Spark application.
type SparkApplicationType string

// Different types of Spark applications.
const (
	JavaApplicationType   SparkApplicationType = "Java"
	ScalaApplicationType  SparkApplicationType = "Scala"
	PythonApplicationType SparkApplicationType = "Python"
	RApplicationType      SparkApplicationType = "R"
)

// DeployMode describes the type of deployment of a Spark application.
type DeployMode string

// Different types of deployments.
const (
	ClusterMode         DeployMode = "cluster"
	ClientMode          DeployMode = "client"
	InClusterClientMode DeployMode = "in-cluster-client"
)

// RestartPolicy is the policy of if and in which conditions the controller should restart a terminated application.
// The decision is based on the policy and termination state of the driver pod. An application that fails submission
// won't be restarted regardless of the policy.
type RestartPolicy string

const (
	Undefined RestartPolicy = ""
	Never     RestartPolicy = "Never"
	OnFailure RestartPolicy = "OnFailure"
	Always    RestartPolicy = "Always"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true

type ScheduledSparkApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ScheduledSparkApplicationSpec   `json:"spec"`
	Status            ScheduledSparkApplicationStatus `json:"status,omitempty"`
}

type ConcurrencyPolicy string

const (
	// ConcurrencyAllow allows SparkApplications to run concurrently.
	ConcurrencyAllow ConcurrencyPolicy = "Allow"
	// ConcurrencyForbid forbids concurrent runs of SparkApplications, skipping the next run if the previous
	// one hasn't finished yet.
	ConcurrencyForbid ConcurrencyPolicy = "Forbid"
	// ConcurrencyReplace kills the currently running SparkApplication instance and replaces it with a new one.
	ConcurrencyReplace ConcurrencyPolicy = "Replace"
)

type ScheduledSparkApplicationSpec struct {
	// Schedule is a cron schedule on which the application should run.
	Schedule string `json:"schedule"`
	// Template is a template from which SparkApplication instances can be created.
	Template SparkApplicationSpec `json:"template"`
	// Suspend is a flag telling the controller to suspend subsequent runs of the application if set to true.
	// Optional.
	// Defaults to false.
	Suspend *bool `json:"suspend,omitempty"`
	// ConcurrencyPolicy is the policy governing concurrent SparkApplication runs.
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`
	// SuccessfulRunHistoryLimit is the number of past successful runs of the application to keep.
	// Optional.
	// Defaults to 1.
	SuccessfulRunHistoryLimit *int32 `json:"successfulRunHistoryLimit,omitempty"`
	// FailedRunHistoryLimit is the number of past failed runs of the application to keep.
	// Optional.
	// Defaults to 1.
	FailedRunHistoryLimit *int32 `json:"failedRunHistoryLimit,omitempty"`
}

type ScheduleState string

const (
	FailedValidationState ScheduleState = "FailedValidation"
	ScheduledState        ScheduleState = "Scheduled"
)

type ScheduledSparkApplicationStatus struct {
	// LastRun is the time when the last run of the application started.
	LastRun metav1.Time `json:"lastRun,omitempty"`
	// NextRun is the time when the next run of the application will start.
	NextRun metav1.Time `json:"nextRun,omitempty"`
	// LastRunName is the name of the SparkApplication for the most recent run of the application.
	LastRunName string `json:"lastRunName,omitempty"`
	// PastSuccessfulRunNames keeps the names of SparkApplications for past successful runs.
	PastSuccessfulRunNames []string `json:"pastSuccessfulRunNames,omitempty"`
	// PastFailedRunNames keeps the names of SparkApplications for past failed runs.
	PastFailedRunNames []string `json:"pastFailedRunNames,omitempty"`
	// ScheduleState is the current scheduling state of the application.
	ScheduleState ScheduleState `json:"scheduleState,omitempty"`
	// Reason tells why the ScheduledSparkApplication is in the particular ScheduleState.
	Reason string `json:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ScheduledSparkApplicationList carries a list of ScheduledSparkApplication objects.
type ScheduledSparkApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduledSparkApplication `json:"items,omitempty"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true

// SparkApplication represents a Spark application running on and using Kubernetes as a cluster manager.
type SparkApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              SparkApplicationSpec   `json:"spec"`
	Status            SparkApplicationStatus `json:"status,omitempty"`
}

// SparkApplicationSpec describes the specification of a Spark application using Kubernetes as a cluster manager.
// It carries every pieces of information a spark-submit command takes and recognizes.
type SparkApplicationSpec struct {
	// Type tells the type of the Spark application.
	Type SparkApplicationType `json:"type"`
	// Mode is the deployment mode of the Spark application.
	Mode DeployMode `json:"mode,omitempty"`
	// Image is the container image for the driver, executor, and init-container. Any custom container images for the
	// driver, executor, or init-container takes precedence over this.
	// Optional.
	Image *string `json:"image,omitempty"`
	// InitContainerImage is the image of the init-container to use. Overrides Spec.Image if set.
	// Optional.
	InitContainerImage *string `json:"initContainerImage,omitempty"`
	// ImagePullPolicy is the image pull policy for the driver, executor, and init-container.
	// Optional.
	ImagePullPolicy *string `json:"imagePullPolicy,omitempty"`
	// ImagePullSecrets is the list of image-pull secrets.
	// Optional.
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
	// MainClass is the fully-qualified main class of the Spark application.
	// This only applies to Java/Scala Spark applications.
	// Optional.
	MainClass *string `json:"mainClass,omitempty"`
	// MainFile is the path to a bundled JAR, Python, or R file of the application.
	// Optional.
	MainApplicationFile *string `json:"mainApplicationFile"`
	// Arguments is a list of arguments to be passed to the application.
	// Optional.
	Arguments []string `json:"arguments,omitempty"`
	// SparkConf carries user-specified Spark configuration properties as they would use the  "--conf" option in
	// spark-submit.
	// Optional.
	SparkConf map[string]string `json:"sparkConf,omitempty"`
	// HadoopConf carries user-specified Hadoop configuration properties as they would use the  the "--conf" option
	// in spark-submit.  The SparkApplication controller automatically adds prefix "spark.hadoop." to Hadoop
	// configuration properties.
	// Optional.
	HadoopConf map[string]string `json:"hadoopConf,omitempty"`
	// SparkConfigMap carries the name of the ConfigMap containing Spark configuration files such as log4j.properties.
	// The controller will add environment variable SPARK_CONF_DIR to the path where the ConfigMap is mounted to.
	// Optional.
	SparkConfigMap *string `json:"sparkConfigMap,omitempty"`
	// HadoopConfigMap carries the name of the ConfigMap containing Hadoop configuration files such as core-site.xml.
	// The controller will add environment variable HADOOP_CONF_DIR to the path where the ConfigMap is mounted to.
	// Optional.
	HadoopConfigMap *string `json:"hadoopConfigMap,omitempty"`
	// Volumes is the list of Kubernetes volumes that can be mounted by the driver and/or executors.
	// Optional.
	Volumes []apiv1.Volume `json:"volumes,omitempty"`
	// Driver is the driver specification.
	Driver DriverSpec `json:"driver"`
	// Executor is the executor specification.
	Executor ExecutorSpec `json:"executor"`
	// Deps captures all possible types of dependencies of a Spark application.
	Deps Dependencies `json:"deps"`
	// RestartPolicy defines the policy on if and in which conditions the controller should restart a failed application.
	RestartPolicy RestartPolicy `json:"restartPolicy,omitempty"`
	// NodeSelector is the Kubernetes node selector to be added to the driver and executor pods.
	// Optional.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// MaxSubmissionRetries is the maximum number of times to retry a failed submission.
	// Optional.
	MaxSubmissionRetries *int32 `json:"maxSubmissionRetries,omitempty"`
	// SubmissionRetryInterval is the unit of intervals in seconds between submission retries.
	// Optional.
	SubmissionRetryInterval *int64 `json:"submissionRetryInterval,omitempty"`
}

// ApplicationStateType represents the type of the current state of an application.
type ApplicationStateType string

// Different states an application may have.
const (
	NewState              ApplicationStateType = "NEW"
	SubmittedState        ApplicationStateType = "SUBMITTED"
	RunningState          ApplicationStateType = "RUNNING"
	CompletedState        ApplicationStateType = "COMPLETED"
	FailedState           ApplicationStateType = "FAILED"
	FailedSubmissionState ApplicationStateType = "SUBMISSION_FAILED"
	UnknownState          ApplicationStateType = "UNKNOWN"
)

// ApplicationState tells the current state of the application and an error message in case of failures.
type ApplicationState struct {
	State        ApplicationStateType `json:"state"`
	ErrorMessage string               `json:"errorMessage"`
}

// ExecutorState tells the current state of an executor.
type ExecutorState string

// Different states an executor may have.
const (
	ExecutorPendingState   ExecutorState = "PENDING"
	ExecutorRunningState   ExecutorState = "RUNNING"
	ExecutorCompletedState ExecutorState = "COMPLETED"
	ExecutorFailedState    ExecutorState = "FAILED"
	ExecutorUnknownState   ExecutorState = "UNKNOWN"
)

// SparkApplicationStatus describes the current status of a Spark application.
type SparkApplicationStatus struct {
	// AppId is the application ID that's also added as a label to the SparkApplication object
	// and driver and executor Pods, and is used to group the objects for the same application.
	AppID string `json:"appId,omitempty"`
	// SubmissionTime is the time when the application is submitted.
	SubmissionTime metav1.Time `json:"submissionTime,omitempty"`
	// CompletionTime is the time when the application runs to completion if it does.
	CompletionTime metav1.Time `json:"completionTime,omitempty"`
	// DriverInfo has information about the driver.
	DriverInfo DriverInfo `json:"driverInfo"`
	// AppState tells the overall application state.
	AppState ApplicationState `json:"applicationState,omitempty"`
	// ExecutorState records the state of executors by executor Pod names.
	ExecutorState map[string]ExecutorState `json:"executorState,omitempty"`
	// SubmissionRetries is the number of retries attempted for a failed submission.
	SubmissionRetries int32 `json:"submissionRetries,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SparkApplicationList carries a list of SparkApplication objects.
type SparkApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SparkApplication `json:"items,omitempty"`
}

// Dependencies specifies all possible types of dependencies of a Spark application.
type Dependencies struct {
	// Jars is a list of JAR files the Spark application depends on.
	// Optional.
	Jars []string `json:"jars,omitempty"`
	// Files is a list of files the Spark application depends on.
	// Optional.
	Files []string `json:"files,omitempty"`
	// PyFiles is a list of Python files the Spark application depends on.
	// Optional.
	PyFiles []string `json:"pyFiles,omitempty"`
	// JarsDownloadDir is the location to download jars to in the driver and executors.
	JarsDownloadDir *string `json:"jarsDownloadDir,omitempty"`
	// FilesDownloadDir is the location to download files to in the driver and executors.
	FilesDownloadDir *string `json:"filesDownloadDir,omitempty"`
	// DownloadTimeout specifies the timeout in seconds before aborting the attempt to download
	// and unpack dependencies from remote locations into the driver and executor pods.
	DownloadTimeout *int32 `json:"downloadTimeout,omitempty"`
	// MaxSimultaneousDownloads specifies the maximum number of remote dependencies to download
	// simultaneously in a driver or executor pod.
	MaxSimultaneousDownloads *int32 `json:"maxSimultaneousDownloads,omitempty"`
}

// SparkPodSpec defines common things that can be customized for a Spark driver or executor pod.
// TODO: investigate if we should use v1.PodSpec and limit what can be set instead.
type SparkPodSpec struct {
	// Cores is the number of CPU cores to request for the pod.
	// Optional.
	Cores *float32 `json:"cores,omitempty"`
	// CoreLimit specifies a hard limit on CPU cores for the pod.
	// Optional
	CoreLimit *string `json:"coreLimit,omitempty"`
	// Memory is the amount of memory to request for the pod.
	// Optional.
	Memory *string `json:"memory,omitempty"`
	// MemoryOverhead is the amount of off-heap memory to allocate in cluster mode, in MiB unless otherwise specified.
	// Optional.
	MemoryOverhead *string `json:"memoryOverhead,omitempty"`
	// Image is the container image to use. Overrides Spec.Image if set.
	// Optional.
	Image *string `json:"image,omitempty"`
	// ConfigMaps carries information of other ConfigMaps to add to the pod.
	// Optional.
	ConfigMaps []NamePath `json:"configMaps,omitempty"`
	// Secrets carries information of secrets to add to the pod.
	// Optional.
	Secrets []SecretInfo `json:"secrets,omitempty"`
	// EnvVars carries the environment variables to add to the pod.
	// Optional.
	EnvVars map[string]string `json:"envVars,omitempty"`
	// EnvSecretKeyRefs holds a mapping from environment variable names to SecretKeyRefs.
	// Optional.
	EnvSecretKeyRefs map[string]NameKey `json:"envSecretKeyRefs,omitempty"`
	// Labels are the Kubernetes labels to be added to the pod.
	// Optional.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are the Kubernetes annotations to be added to the pod.
	// Optional.
	Annotations map[string]string `json:"annotations,omitempty"`
	// VolumeMounts specifies the volumes listed in ".spec.volumes" to mount into the main container's filesystem.
	// Optional.
	VolumeMounts []apiv1.VolumeMount `json:"volumeMounts,omitempty"`
	// Affinity specifies the affinity/anti-affinity settings for the pod.
	// Optional.
	Affinity *apiv1.Affinity `json:"affinity,omitempty"`
}

// DriverSpec is specification of the driver.
type DriverSpec struct {
	SparkPodSpec
	// PodName is the name of the driver pod that the user creates. This is used for the
	// in-cluster client mode in which the user creates a client pod where the driver of
	// the user application runs. It's an error to set this field if Mode is not
	// in-cluster-client.
	// Optional.
	PodName *string `json:"podName,omitempty"`
	// ServiceAccount is the name of the Kubernetes service account used by the driver pod
	// when requesting executor pods from the API server.
	ServiceAccount *string `json:"serviceAccount,omitempty"`
}

// ExecutorSpec is specification of the executor.
type ExecutorSpec struct {
	SparkPodSpec
	// Instances is the number of executor instances.
	// Optional.
	Instances *int32 `json:"instances,omitempty"`
	// CoreRequest is the physical CPU core request for the executors.
	// Optional.
	CoreRequest *string `json:"coreRequest,omitempty"`
}

// NamePath is a pair of a name and a path to which the named objects should be mounted to.
type NamePath struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// SecretType tells the type of a secret.
type SecretType string

// An enumeration of secret types supported.
const (
	// GCPServiceAccountSecret is for secrets from a GCP service account Json key file that needs
	// the environment variable GOOGLE_APPLICATION_CREDENTIALS.
	GCPServiceAccountSecret SecretType = "GCPServiceAccount"
	// HadoopDelegationTokenSecret is for secrets from an Hadoop delegation token that needs the
	// environment variable HADOOP_TOKEN_FILE_LOCATION.
	HadoopDelegationTokenSecret SecretType = "HadoopDelegationToken"
	// GenericType is for secrets that needs no special handling.
	GenericType SecretType = "Generic"
)

// DriverInfo captures information about the driver.
type DriverInfo struct {
	WebUIServiceName string `json:"webUIServiceName,omitempty"`
	WebUIPort        int32  `json:"webUIPort,omitempty"`
	WebUIAddress     string `json:"webUIAddress,omitempty"`
	PodName          string `json:"podName,omitempty"`
}

// SecretInfo captures information of a secret.
type SecretInfo struct {
	Name string     `json:"name"`
	Path string     `json:"path"`
	Type SecretType `json:"secretType"`
}

// NameKey represents the name and key of a SecretKeyRef.
type NameKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
