// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	"errors"
	"fmt"

	util "github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/operator-tools/pkg/volume"
	"github.com/spf13/cast"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +name:"LoggingSpec"
// +weight:"200"
type _hugoLoggingSpec interface{}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +name:"Logging"
// +version:"v1beta1"
// +description:"Logging system configuration"
type _metaLoggingSpec interface{}

// LoggingSpec defines the desired state of Logging
type LoggingSpec struct {
	LoggingRef                             string           `json:"loggingRef,omitempty"`
	FlowConfigCheckDisabled                bool             `json:"flowConfigCheckDisabled,omitempty"`
	FlowConfigOverride                     string           `json:"flowConfigOverride,omitempty"`
	FluentbitSpec                          *FluentbitSpec   `json:"fluentbit,omitempty"`
	FluentdSpec                            *FluentdSpec     `json:"fluentd,omitempty"`
	DefaultFlowSpec                        *DefaultFlowSpec `json:"defaultFlow,omitempty"`
	GlobalFilters                          []Filter         `json:"globalFilters,omitempty"`
	WatchNamespaces                        []string         `json:"watchNamespaces,omitempty"`
	ControlNamespace                       string           `json:"controlNamespace"`
	AllowClusterResourcesFromAllNamespaces bool             `json:"allowClusterResourcesFromAllNamespaces,omitempty"`

	// EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the
	// fluentbit daemonset and the fluentd statefulset (and possibly other resource in the future)
	// in case there is a change in an immutable field
	// that otherwise couldn't be managed with a simple update.
	EnableRecreateWorkloadOnImmutableFieldChange bool `json:"enableRecreateWorkloadOnImmutableFieldChange,omitempty"`
}

// LoggingStatus defines the observed state of Logging
type LoggingStatus struct {
	ConfigCheckResults map[string]bool `json:"configCheckResults,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=loggings,scope=Cluster,categories=logging-all

// Logging is the Schema for the loggings API
type Logging struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoggingSpec   `json:"spec,omitempty"`
	Status LoggingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoggingList contains a list of Logging
type LoggingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Logging `json:"items"`
}

// +kubebuilder:object:generate=true

// DefaultFlowSpec is a Flow for logs that did not match any other Flow
type DefaultFlowSpec struct {
	Filters []Filter `json:"filters,omitempty"`
	// Deprecated
	OutputRefs       []string `json:"outputRefs,omitempty"`
	GlobalOutputRefs []string `json:"globalOutputRefs,omitempty"`
}

const (
	DefaultFluentbitImageRepository = "fluent/fluent-bit"
	DefaultFluentbitImageTag        = "1.6.4"
	DefaultFluentdImageRepository   = "ghcr.io/banzaicloud/fluentd"
	DefaultFluentdImageTag          = "v1.11.5-alpine-1"
)

// SetDefaults fills empty attributes
func (l *Logging) SetDefaults() error {
	if !l.Spec.FlowConfigCheckDisabled && l.Status.ConfigCheckResults == nil {
		l.Status.ConfigCheckResults = make(map[string]bool)
	}
	if l.Spec.FluentdSpec != nil {
		if l.Spec.FluentdSpec.FluentdPvcSpec != nil {
			return errors.New("`fluentdPvcSpec` field is deprecated, use: `bufferStorageVolume`")
		}
		if l.Spec.FluentdSpec.Image.Repository == "" {
			l.Spec.FluentdSpec.Image.Repository = DefaultFluentdImageRepository
		}
		if l.Spec.FluentdSpec.Image.Tag == "" {
			l.Spec.FluentdSpec.Image.Tag = DefaultFluentdImageTag
		}
		if l.Spec.FluentdSpec.Image.PullPolicy == "" {
			l.Spec.FluentdSpec.Image.PullPolicy = "IfNotPresent"
		}
		if l.Spec.FluentdSpec.Annotations == nil {
			l.Spec.FluentdSpec.Annotations = make(map[string]string)
		}
		if l.Spec.FluentdSpec.Security == nil {
			l.Spec.FluentdSpec.Security = &Security{}
		}
		if l.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate == nil {
			l.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate = util.BoolPointer(true)
		}
		if l.Spec.FluentdSpec.Security.SecurityContext == nil {
			l.Spec.FluentdSpec.Security.SecurityContext = &v1.SecurityContext{}
		}
		if l.Spec.FluentdSpec.Security.PodSecurityContext == nil {
			l.Spec.FluentdSpec.Security.PodSecurityContext = &v1.PodSecurityContext{}
		}
		if l.Spec.FluentdSpec.Security.PodSecurityContext.FSGroup == nil {
			l.Spec.FluentdSpec.Security.PodSecurityContext.FSGroup = util.IntPointer64(101)
		}
		if l.Spec.FluentdSpec.Metrics != nil {
			if l.Spec.FluentdSpec.Metrics.Path == "" {
				l.Spec.FluentdSpec.Metrics.Path = "/metrics"
			}
			if l.Spec.FluentdSpec.Metrics.Port == 0 {
				l.Spec.FluentdSpec.Metrics.Port = 24231
			}
			if l.Spec.FluentdSpec.Metrics.Timeout == "" {
				l.Spec.FluentdSpec.Metrics.Timeout = "5s"
			}
			if l.Spec.FluentdSpec.Metrics.Interval == "" {
				l.Spec.FluentdSpec.Metrics.Interval = "15s"
			}

			if l.Spec.FluentdSpec.Metrics.PrometheusAnnotations {
				l.Spec.FluentdSpec.Annotations["prometheus.io/scrape"] = "true"

				l.Spec.FluentdSpec.Annotations["prometheus.io/path"] = l.Spec.FluentdSpec.Metrics.Path
				l.Spec.FluentdSpec.Annotations["prometheus.io/port"] = fmt.Sprintf("%d", l.Spec.FluentdSpec.Metrics.Port)
			}
		}

		if !l.Spec.FluentdSpec.DisablePvc {
			if l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim == nil {
				l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim = &volume.PersistentVolumeClaim{
					PersistentVolumeClaimSpec: v1.PersistentVolumeClaimSpec{},
				}
			}
			if l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec.AccessModes == nil {
				l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec.AccessModes = []v1.PersistentVolumeAccessMode{
					v1.ReadWriteOnce,
				}
			}
			if l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec.Resources.Requests == nil {
				l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec.Resources.Requests = map[v1.ResourceName]resource.Quantity{
					"storage": resource.MustParse("20Gi"),
				}
			}
			if l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec.VolumeMode == nil {
				l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec.VolumeMode = persistentVolumeModePointer(v1.PersistentVolumeFilesystem)
			}
			if l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName == "" {
				l.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName = "fluentd-buffer"
			}
		}
		if l.Spec.FluentdSpec.VolumeModImage.Repository == "" {
			l.Spec.FluentdSpec.VolumeModImage.Repository = "busybox"
		}
		if l.Spec.FluentdSpec.VolumeModImage.Tag == "" {
			l.Spec.FluentdSpec.VolumeModImage.Tag = "latest"
		}
		if l.Spec.FluentdSpec.VolumeModImage.PullPolicy == "" {
			l.Spec.FluentdSpec.VolumeModImage.PullPolicy = "IfNotPresent"
		}
		if l.Spec.FluentdSpec.ConfigReloaderImage.Repository == "" {
			l.Spec.FluentdSpec.ConfigReloaderImage.Repository = "jimmidyson/configmap-reload"
		}
		if l.Spec.FluentdSpec.ConfigReloaderImage.Tag == "" {
			l.Spec.FluentdSpec.ConfigReloaderImage.Tag = "v0.2.2"
		}
		if l.Spec.FluentdSpec.ConfigReloaderImage.PullPolicy == "" {
			l.Spec.FluentdSpec.ConfigReloaderImage.PullPolicy = "IfNotPresent"
		}
		if l.Spec.FluentdSpec.Resources.Limits == nil {
			l.Spec.FluentdSpec.Resources.Limits = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("400M"),
				v1.ResourceCPU:    resource.MustParse("1000m"),
			}
		}
		if l.Spec.FluentdSpec.Resources.Requests == nil {
			l.Spec.FluentdSpec.Resources.Requests = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("100M"),
				v1.ResourceCPU:    resource.MustParse("500m"),
			}
		}
		if l.Spec.FluentdSpec.Port == 0 {
			l.Spec.FluentdSpec.Port = 24240
		}
		if l.Spec.FluentdSpec.Scaling == nil {
			l.Spec.FluentdSpec.Scaling = new(FluentdScaling)
		}
		if l.Spec.FluentdSpec.Scaling.Replicas == 0 {
			l.Spec.FluentdSpec.Scaling.Replicas = 1
		}
		if l.Spec.FluentdSpec.Scaling.PodManagementPolicy == "" {
			l.Spec.FluentdSpec.Scaling.PodManagementPolicy = "OrderedReady"
		}
		if l.Spec.FluentdSpec.FluentLogDestination == "" {
			l.Spec.FluentdSpec.FluentLogDestination = "null"
		}
		if l.Spec.FluentdSpec.FluentOutLogrotate == nil {
			l.Spec.FluentdSpec.FluentOutLogrotate = &FluentOutLogrotate{
				Enabled: true,
			}
		}
		if l.Spec.FluentdSpec.FluentOutLogrotate.Path == "" {
			l.Spec.FluentdSpec.FluentOutLogrotate.Path = "/fluentd/log/out"
		}
		if l.Spec.FluentdSpec.FluentOutLogrotate.Age == "" {
			l.Spec.FluentdSpec.FluentOutLogrotate.Age = "10"
		}
		if l.Spec.FluentdSpec.FluentOutLogrotate.Size == "" {
			l.Spec.FluentdSpec.FluentOutLogrotate.Size = cast.ToString(1024 * 1024 * 10)
		}
		if l.Spec.FluentdSpec.LivenessProbe == nil {
			if l.Spec.FluentdSpec.LivenessDefaultCheck {
				l.Spec.FluentdSpec.LivenessProbe = &v1.Probe{
					Handler: v1.Handler{
						Exec: &v1.ExecAction{Command: []string{"/bin/healthy.sh"}},
					},
					InitialDelaySeconds: 600,
					TimeoutSeconds:      0,
					PeriodSeconds:       60,
					SuccessThreshold:    0,
					FailureThreshold:    0,
				}
			}
		}
	}
	if l.Spec.FluentbitSpec != nil {
		if l.Spec.FluentbitSpec.PosisionDBLegacy != nil {
			return errors.New("`position_db` field is deprecated, use `positiondb`")
		}
		if l.Spec.FluentbitSpec.Parser != "" {
			return errors.New("`parser` field is deprecated, use `inputTail.Parser`")
		}
		if l.Spec.FluentbitSpec.Image.Repository == "" {
			l.Spec.FluentbitSpec.Image.Repository = DefaultFluentbitImageRepository
		}
		if l.Spec.FluentbitSpec.Image.Tag == "" {
			l.Spec.FluentbitSpec.Image.Tag = DefaultFluentbitImageTag
		}
		if l.Spec.FluentbitSpec.Image.PullPolicy == "" {
			l.Spec.FluentbitSpec.Image.PullPolicy = "IfNotPresent"
		}
		if l.Spec.FluentbitSpec.Flush == 0 {
			l.Spec.FluentbitSpec.Flush = 1
		}
		if l.Spec.FluentbitSpec.Grace == 0 {
			l.Spec.FluentbitSpec.Grace = 5
		}
		if l.Spec.FluentbitSpec.LogLevel == "" {
			l.Spec.FluentbitSpec.LogLevel = "info"
		}
		if l.Spec.FluentbitSpec.CoroStackSize == 0 {
			l.Spec.FluentbitSpec.CoroStackSize = 24576
		}
		if l.Spec.FluentbitSpec.Resources.Limits == nil {
			l.Spec.FluentbitSpec.Resources.Limits = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("100M"),
				v1.ResourceCPU:    resource.MustParse("200m"),
			}
		}
		if l.Spec.FluentbitSpec.Resources.Requests == nil {
			l.Spec.FluentbitSpec.Resources.Requests = v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("50M"),
				v1.ResourceCPU:    resource.MustParse("100m"),
			}
		}
		if l.Spec.FluentbitSpec.InputTail.Path == "" {
			l.Spec.FluentbitSpec.InputTail.Path = "/var/log/containers/*.log"
		}
		if l.Spec.FluentbitSpec.InputTail.RefreshInterval == "" {
			l.Spec.FluentbitSpec.InputTail.RefreshInterval = "5"
		}
		if l.Spec.FluentbitSpec.InputTail.SkipLongLines == "" {
			l.Spec.FluentbitSpec.InputTail.SkipLongLines = "On"
		}
		if l.Spec.FluentbitSpec.InputTail.DB == nil {
			l.Spec.FluentbitSpec.InputTail.DB = util.StringPointer("/tail-db/tail-containers-state.db")
		}
		if l.Spec.FluentbitSpec.InputTail.MemBufLimit == "" {
			l.Spec.FluentbitSpec.InputTail.MemBufLimit = "5MB"
		}
		if l.Spec.FluentbitSpec.InputTail.Tag == "" {
			l.Spec.FluentbitSpec.InputTail.Tag = "kubernetes.*"
		}
		if l.Spec.FluentbitSpec.Annotations == nil {
			l.Spec.FluentbitSpec.Annotations = make(map[string]string)
		}
		if l.Spec.FluentbitSpec.Security == nil {
			l.Spec.FluentbitSpec.Security = &Security{}
		}
		if l.Spec.FluentbitSpec.Security.RoleBasedAccessControlCreate == nil {
			l.Spec.FluentbitSpec.Security.RoleBasedAccessControlCreate = util.BoolPointer(true)
		}
		if l.Spec.FluentbitSpec.Security.SecurityContext == nil {
			l.Spec.FluentbitSpec.Security.SecurityContext = &v1.SecurityContext{}
		}
		if l.Spec.FluentbitSpec.Security.PodSecurityContext == nil {
			l.Spec.FluentbitSpec.Security.PodSecurityContext = &v1.PodSecurityContext{}
		}
		if l.Spec.FluentbitSpec.Metrics != nil {
			if l.Spec.FluentbitSpec.Metrics.Path == "" {
				l.Spec.FluentbitSpec.Metrics.Path = "/api/v1/metrics/prometheus"
			}
			if l.Spec.FluentbitSpec.Metrics.Port == 0 {
				l.Spec.FluentbitSpec.Metrics.Port = 2020
			}
			if l.Spec.FluentbitSpec.Metrics.Timeout == "" {
				l.Spec.FluentbitSpec.Metrics.Timeout = "5s"
			}
			if l.Spec.FluentbitSpec.Metrics.Interval == "" {
				l.Spec.FluentbitSpec.Metrics.Interval = "15s"
			}
			if l.Spec.FluentbitSpec.Metrics.PrometheusAnnotations {
				l.Spec.FluentbitSpec.Annotations["prometheus.io/scrape"] = "true"
				l.Spec.FluentbitSpec.Annotations["prometheus.io/path"] = l.Spec.FluentbitSpec.Metrics.Path
				l.Spec.FluentbitSpec.Annotations["prometheus.io/port"] = fmt.Sprintf("%d", l.Spec.FluentbitSpec.Metrics.Port)
			}
		} else if l.Spec.FluentbitSpec.LivenessDefaultCheck {
			l.Spec.FluentbitSpec.Metrics = &Metrics{
				Port: 2020,
				Path: "/",
			}
		}
		if l.Spec.FluentbitSpec.LivenessProbe == nil {
			if l.Spec.FluentbitSpec.LivenessDefaultCheck {
				l.Spec.FluentbitSpec.LivenessProbe = &v1.Probe{
					Handler: v1.Handler{
						HTTPGet: &v1.HTTPGetAction{
							Path: l.Spec.FluentbitSpec.Metrics.Path,
							Port: intstr.IntOrString{
								IntVal: l.Spec.FluentbitSpec.Metrics.Port,
							},
						}},
					InitialDelaySeconds: 10,
					TimeoutSeconds:      0,
					PeriodSeconds:       10,
					SuccessThreshold:    0,
					FailureThreshold:    3,
				}
			}
		}

		if l.Spec.FluentbitSpec.MountPath == "" {
			l.Spec.FluentbitSpec.MountPath = "/var/lib/docker/containers"
		}
		if l.Spec.FluentbitSpec.BufferStorage.StoragePath == "" {
			l.Spec.FluentbitSpec.BufferStorage.StoragePath = "/buffers"
		}
		if l.Spec.FluentbitSpec.FilterAws != nil {
			if l.Spec.FluentbitSpec.FilterAws.ImdsVersion == "" {
				l.Spec.FluentbitSpec.FilterAws.ImdsVersion = "v2"
			}
			if l.Spec.FluentbitSpec.FilterAws.AZ == nil {
				l.Spec.FluentbitSpec.FilterAws.AZ = util.BoolPointer(true)
			}
			if l.Spec.FluentbitSpec.FilterAws.Ec2InstanceID == nil {
				l.Spec.FluentbitSpec.FilterAws.Ec2InstanceID = util.BoolPointer(true)
			}
			if l.Spec.FluentbitSpec.FilterAws.Ec2InstanceType == nil {
				l.Spec.FluentbitSpec.FilterAws.Ec2InstanceType = util.BoolPointer(false)
			}
			if l.Spec.FluentbitSpec.FilterAws.PrivateIP == nil {
				l.Spec.FluentbitSpec.FilterAws.PrivateIP = util.BoolPointer(false)
			}
			if l.Spec.FluentbitSpec.FilterAws.AmiID == nil {
				l.Spec.FluentbitSpec.FilterAws.AmiID = util.BoolPointer(false)
			}
			if l.Spec.FluentbitSpec.FilterAws.AccountID == nil {
				l.Spec.FluentbitSpec.FilterAws.AccountID = util.BoolPointer(false)
			}
			if l.Spec.FluentbitSpec.FilterAws.Hostname == nil {
				l.Spec.FluentbitSpec.FilterAws.Hostname = util.BoolPointer(false)
			}
			if l.Spec.FluentbitSpec.FilterAws.VpcID == nil {
				l.Spec.FluentbitSpec.FilterAws.VpcID = util.BoolPointer(false)
			}
		}
		if l.Spec.FluentbitSpec.ForwardOptions == nil {
			l.Spec.FluentbitSpec.ForwardOptions = &ForwardOptions{}
		}
		if l.Spec.FluentbitSpec.ForwardOptions.RetryLimit == "" {
			l.Spec.FluentbitSpec.ForwardOptions.RetryLimit = "False"
		}
	}
	return nil
}

// SetDefaultsOnCopy makes a deep copy of the instance and sets defaults on the copy
func (l *Logging) SetDefaultsOnCopy() (*Logging, error) {
	if l == nil {
		return nil, nil
	}

	copy := l.DeepCopy()
	if err := copy.SetDefaults(); err != nil {
		return nil, err
	}
	return copy, nil
}

// QualifiedName is the "logging-resource" name combined
func (l *Logging) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", l.Name, name)
}

func init() {
	SchemeBuilder.Register(&Logging{}, &LoggingList{})
}

func persistentVolumeModePointer(mode v1.PersistentVolumeMode) *v1.PersistentVolumeMode {
	return &mode
}
