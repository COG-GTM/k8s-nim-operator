/*
Copyright 2024.

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
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/NVIDIA/k8s-nim-operator/api/apps/v1beta1"
)

// ConvertTo converts this NIMCache (v1alpha1) to the Hub version (v1beta1).
func (src *NIMCache) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.NIMCache)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec fields
	dst.Spec.Source = convertNIMSourceToV1beta1(src.Spec.Source)
	dst.Spec.Storage = convertNIMCacheStorageToV1beta1(src.Spec.Storage)
	dst.Spec.Resources = convertResourcesToV1beta1(src.Spec.Resources)
	dst.Spec.Tolerations = src.Spec.Tolerations
	dst.Spec.NodeSelector = src.Spec.NodeSelector
	dst.Spec.UserID = src.Spec.UserID
	dst.Spec.GroupID = src.Spec.GroupID
	dst.Spec.Env = src.Spec.Env
	dst.Spec.RuntimeClassName = src.Spec.RuntimeClassName

	// Handle Proxy conversion
	// If Proxy is set, use it directly
	// If CertConfig is set but Proxy is not, create a Proxy with the CertConfig's ConfigMap name
	if src.Spec.Proxy != nil {
		dst.Spec.Proxy = convertProxySpecToV1beta1(src.Spec.Proxy)
	} else if src.Spec.CertConfig != nil {
		// Convert deprecated CertConfig to Proxy
		dst.Spec.Proxy = &v1beta1.ProxySpec{
			CertConfigMap: src.Spec.CertConfig.Name,
		}
	}

	// Status fields
	dst.Status.State = src.Status.State
	dst.Status.PVC = src.Status.PVC
	dst.Status.Profiles = convertNIMProfilesToV1beta1(src.Status.Profiles)
	dst.Status.Conditions = src.Status.Conditions

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this version (v1alpha1).
func (dst *NIMCache) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.NIMCache)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Spec fields
	dst.Spec.Source = convertNIMSourceFromV1beta1(src.Spec.Source)
	dst.Spec.Storage = convertNIMCacheStorageFromV1beta1(src.Spec.Storage)
	dst.Spec.Resources = convertResourcesFromV1beta1(src.Spec.Resources)
	dst.Spec.Tolerations = src.Spec.Tolerations
	dst.Spec.NodeSelector = src.Spec.NodeSelector
	dst.Spec.UserID = src.Spec.UserID
	dst.Spec.GroupID = src.Spec.GroupID
	dst.Spec.Env = src.Spec.Env
	dst.Spec.RuntimeClassName = src.Spec.RuntimeClassName

	// Handle Proxy conversion
	if src.Spec.Proxy != nil {
		dst.Spec.Proxy = convertProxySpecFromV1beta1(src.Spec.Proxy)
		// Also populate deprecated CertConfig for backward compatibility if CertConfigMap is set
		if src.Spec.Proxy.CertConfigMap != "" {
			dst.Spec.CertConfig = &CertConfig{
				Name:      src.Spec.Proxy.CertConfigMap,
				MountPath: "/etc/ssl/certs", // Default mount path
			}
		}
	}

	// Status fields
	dst.Status.State = src.Status.State
	dst.Status.PVC = src.Status.PVC
	dst.Status.Profiles = convertNIMProfilesFromV1beta1(src.Status.Profiles)
	dst.Status.Conditions = src.Status.Conditions

	return nil
}

// Helper functions for converting NIMSource
func convertNIMSourceToV1beta1(src NIMSource) v1beta1.NIMSource {
	dst := v1beta1.NIMSource{}
	if src.NGC != nil {
		dst.NGC = convertNGCSourceToV1beta1(src.NGC)
	}
	if src.DataStore != nil {
		dst.DataStore = convertNemoDataStoreSourceToV1beta1(src.DataStore)
	}
	if src.HF != nil {
		dst.HF = convertHuggingFaceHubSourceToV1beta1(src.HF)
	}
	return dst
}

func convertNIMSourceFromV1beta1(src v1beta1.NIMSource) NIMSource {
	dst := NIMSource{}
	if src.NGC != nil {
		dst.NGC = convertNGCSourceFromV1beta1(src.NGC)
	}
	if src.DataStore != nil {
		dst.DataStore = convertNemoDataStoreSourceFromV1beta1(src.DataStore)
	}
	if src.HF != nil {
		dst.HF = convertHuggingFaceHubSourceFromV1beta1(src.HF)
	}
	return dst
}

// Helper functions for converting NGCSource
func convertNGCSourceToV1beta1(src *NGCSource) *v1beta1.NGCSource {
	if src == nil {
		return nil
	}
	dst := &v1beta1.NGCSource{
		AuthSecret:    src.AuthSecret,
		ModelPuller:   src.ModelPuller,
		PullSecret:    src.PullSecret,
		ModelEndpoint: src.ModelEndpoint,
	}
	if src.Model != nil {
		dst.Model = convertModelSpecToV1beta1(src.Model)
	}
	return dst
}

func convertNGCSourceFromV1beta1(src *v1beta1.NGCSource) *NGCSource {
	if src == nil {
		return nil
	}
	dst := &NGCSource{
		AuthSecret:    src.AuthSecret,
		ModelPuller:   src.ModelPuller,
		PullSecret:    src.PullSecret,
		ModelEndpoint: src.ModelEndpoint,
	}
	if src.Model != nil {
		dst.Model = convertModelSpecFromV1beta1(src.Model)
	}
	return dst
}

// Helper functions for converting ModelSpec
func convertModelSpecToV1beta1(src *ModelSpec) *v1beta1.ModelSpec {
	if src == nil {
		return nil
	}
	return &v1beta1.ModelSpec{
		Profiles:          src.Profiles,
		Precision:         src.Precision,
		Engine:            src.Engine,
		TensorParallelism: src.TensorParallelism,
		QoSProfile:        src.QoSProfile,
		GPUs:              convertGPUSpecsToV1beta1(src.GPUs),
		Lora:              src.Lora,
		Buildable:         src.Buildable,
	}
}

func convertModelSpecFromV1beta1(src *v1beta1.ModelSpec) *ModelSpec {
	if src == nil {
		return nil
	}
	return &ModelSpec{
		Profiles:          src.Profiles,
		Precision:         src.Precision,
		Engine:            src.Engine,
		TensorParallelism: src.TensorParallelism,
		QoSProfile:        src.QoSProfile,
		GPUs:              convertGPUSpecsFromV1beta1(src.GPUs),
		Lora:              src.Lora,
		Buildable:         src.Buildable,
	}
}

// Helper functions for converting GPUSpec
func convertGPUSpecsToV1beta1(src []GPUSpec) []v1beta1.GPUSpec {
	if src == nil {
		return nil
	}
	dst := make([]v1beta1.GPUSpec, len(src))
	for i, g := range src {
		dst[i] = v1beta1.GPUSpec{
			Product: g.Product,
			IDs:     g.IDs,
		}
	}
	return dst
}

func convertGPUSpecsFromV1beta1(src []v1beta1.GPUSpec) []GPUSpec {
	if src == nil {
		return nil
	}
	dst := make([]GPUSpec, len(src))
	for i, g := range src {
		dst[i] = GPUSpec{
			Product: g.Product,
			IDs:     g.IDs,
		}
	}
	return dst
}

// Helper functions for converting NemoDataStoreSource
func convertNemoDataStoreSourceToV1beta1(src *NemoDataStoreSource) *v1beta1.NemoDataStoreSource {
	if src == nil {
		return nil
	}
	return &v1beta1.NemoDataStoreSource{
		Endpoint:  src.Endpoint,
		Namespace: src.Namespace,
		DSHFCommonFields: v1beta1.DSHFCommonFields{
			ModelName:   src.ModelName,
			DatasetName: src.DatasetName,
			AuthSecret:  src.AuthSecret,
			ModelPuller: src.ModelPuller,
			PullSecret:  src.PullSecret,
			Revision:    src.Revision,
		},
	}
}

func convertNemoDataStoreSourceFromV1beta1(src *v1beta1.NemoDataStoreSource) *NemoDataStoreSource {
	if src == nil {
		return nil
	}
	return &NemoDataStoreSource{
		Endpoint:  src.Endpoint,
		Namespace: src.Namespace,
		DSHFCommonFields: DSHFCommonFields{
			ModelName:   src.ModelName,
			DatasetName: src.DatasetName,
			AuthSecret:  src.AuthSecret,
			ModelPuller: src.ModelPuller,
			PullSecret:  src.PullSecret,
			Revision:    src.Revision,
		},
	}
}

// Helper functions for converting HuggingFaceHubSource
func convertHuggingFaceHubSourceToV1beta1(src *HuggingFaceHubSource) *v1beta1.HuggingFaceHubSource {
	if src == nil {
		return nil
	}
	return &v1beta1.HuggingFaceHubSource{
		Endpoint:  src.Endpoint,
		Namespace: src.Namespace,
		DSHFCommonFields: v1beta1.DSHFCommonFields{
			ModelName:   src.ModelName,
			DatasetName: src.DatasetName,
			AuthSecret:  src.AuthSecret,
			ModelPuller: src.ModelPuller,
			PullSecret:  src.PullSecret,
			Revision:    src.Revision,
		},
	}
}

func convertHuggingFaceHubSourceFromV1beta1(src *v1beta1.HuggingFaceHubSource) *HuggingFaceHubSource {
	if src == nil {
		return nil
	}
	return &HuggingFaceHubSource{
		Endpoint:  src.Endpoint,
		Namespace: src.Namespace,
		DSHFCommonFields: DSHFCommonFields{
			ModelName:   src.ModelName,
			DatasetName: src.DatasetName,
			AuthSecret:  src.AuthSecret,
			ModelPuller: src.ModelPuller,
			PullSecret:  src.PullSecret,
			Revision:    src.Revision,
		},
	}
}

// Helper functions for converting NIMCacheStorage
func convertNIMCacheStorageToV1beta1(src NIMCacheStorage) v1beta1.NIMCacheStorage {
	return v1beta1.NIMCacheStorage{
		PVC: convertPVCToV1beta1(src.PVC),
	}
	// Note: HostPath is deprecated and not converted to v1beta1
}

func convertNIMCacheStorageFromV1beta1(src v1beta1.NIMCacheStorage) NIMCacheStorage {
	return NIMCacheStorage{
		PVC: convertPVCFromV1beta1(src.PVC),
		// HostPath is not populated from v1beta1 as it's deprecated
	}
}

// Helper functions for converting PersistentVolumeClaim
func convertPVCToV1beta1(src PersistentVolumeClaim) v1beta1.PersistentVolumeClaim {
	return v1beta1.PersistentVolumeClaim{
		Create:           src.Create,
		Name:             src.Name,
		StorageClass:     src.StorageClass,
		Size:             src.Size,
		VolumeAccessMode: src.VolumeAccessMode,
		SubPath:          src.SubPath,
		Annotations:      src.Annotations,
	}
}

func convertPVCFromV1beta1(src v1beta1.PersistentVolumeClaim) PersistentVolumeClaim {
	return PersistentVolumeClaim{
		Create:           src.Create,
		Name:             src.Name,
		StorageClass:     src.StorageClass,
		Size:             src.Size,
		VolumeAccessMode: src.VolumeAccessMode,
		SubPath:          src.SubPath,
		Annotations:      src.Annotations,
	}
}

// Helper functions for converting Resources
func convertResourcesToV1beta1(src Resources) v1beta1.Resources {
	return v1beta1.Resources{
		CPU:    src.CPU,
		Memory: src.Memory,
	}
}

func convertResourcesFromV1beta1(src v1beta1.Resources) Resources {
	return Resources{
		CPU:    src.CPU,
		Memory: src.Memory,
	}
}

// Helper functions for converting ProxySpec
func convertProxySpecToV1beta1(src *ProxySpec) *v1beta1.ProxySpec {
	if src == nil {
		return nil
	}
	return &v1beta1.ProxySpec{
		HttpProxy:     src.HttpProxy,
		HttpsProxy:    src.HttpsProxy,
		NoProxy:       src.NoProxy,
		CertConfigMap: src.CertConfigMap,
	}
}

func convertProxySpecFromV1beta1(src *v1beta1.ProxySpec) *ProxySpec {
	if src == nil {
		return nil
	}
	return &ProxySpec{
		HttpProxy:     src.HttpProxy,
		HttpsProxy:    src.HttpsProxy,
		NoProxy:       src.NoProxy,
		CertConfigMap: src.CertConfigMap,
	}
}

// Helper functions for converting NIMProfile
func convertNIMProfilesToV1beta1(src []NIMProfile) []v1beta1.NIMProfile {
	if src == nil {
		return nil
	}
	dst := make([]v1beta1.NIMProfile, len(src))
	for i, p := range src {
		dst[i] = v1beta1.NIMProfile{
			Name:    p.Name,
			Model:   p.Model,
			Release: p.Release,
			Config:  p.Config,
		}
	}
	return dst
}

func convertNIMProfilesFromV1beta1(src []v1beta1.NIMProfile) []NIMProfile {
	if src == nil {
		return nil
	}
	dst := make([]NIMProfile, len(src))
	for i, p := range src {
		dst[i] = NIMProfile{
			Name:    p.Name,
			Model:   p.Model,
			Release: p.Release,
			Config:  p.Config,
		}
	}
	return dst
}
