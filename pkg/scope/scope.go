package scope

import (
	"fmt"

	xclusterv1 "github.com/aizhvaly/machinemetadata/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
)

type MetadataScope struct {
	logr.Logger
	MachineMD          *xclusterv1.MachineMetadata
	WorkloadKubeconfig []byte
}

type MetadataScopeParams struct {
	Logger    logr.Logger
	MachineMD *xclusterv1.MachineMetadata
}

func NewMetadataScope(params *MetadataScopeParams) (*MetadataScope, error) {
	if params.Logger == nil {
		return nil, fmt.Errorf("logger required for MetadataScope params")
	}
	if params.MachineMD == nil {
		return nil, fmt.Errorf("machineMD required for MetadataScope params")
	}

	return &MetadataScope{
		Logger:    params.Logger,
		MachineMD: params.MachineMD,
	}, nil
}

func (m *MetadataScope) SetClusterName(name string) {
	m.MachineMD.ClusterName = name
}

func (m *MetadataScope) GetClusterName() string {
	return m.MachineMD.ClusterName
}

func (m *MetadataScope) SetWorklaodKubeconfig(s *corev1.Secret) error {
	data, ok := s.Data["value"]
	if !ok {
		return fmt.Errorf("default store for capi kubeconfig from secret path 'value' does not exist")
	}

	m.WorkloadKubeconfig = data
	return nil
}

func (m *MetadataScope) GetWorloadKubeconfig() []byte {
	return m.WorkloadKubeconfig
}
