package machinemetadata

import (
	"fmt"

	"k8s.io/client-go/kubernetes/scheme"

	"github.com/aizhvaly/machinemetadata/pkg/scope"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MMD struct {
	client rest.Interface
	scope  *scope.MetadataScope
}

func NewMachineMetadata(scope *scope.MetadataScope) (*MMD, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(scope.GetWorloadKubeconfig())
	if err != nil {
		return nil, fmt.Errorf("failed init kubeconfig, %v", err)
	}

	config.ContentConfig.GroupVersion = &corev1.SchemeGroupVersion
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	cl, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, fmt.Errorf("failed init rest kube client for workload cluster: %v", err)
	}

	return &MMD{
		client: cl,
		scope:  scope,
	}, nil
}

func (m *MMD) Reconcile() error {
	return m.reconcile()
}

func (m *MMD) reconcile() error {
	m.scope.Info("start fetchig worload nodes from targets", "targets", m.scope.MachineMD.Status.Targets)
	var nodes []*corev1.Node
	for _, t := range m.scope.MachineMD.Status.Targets {
		n, err := m.getNode(t, metav1.GetOptions{})
		if err != nil {
			m.scope.Error(err, "failed get node from workload cluster", "node", t)
			continue
		}

		nodes = append(nodes, n)
	}

	for _, node := range nodes {
		m.scope.Info("successful fetch node", "node", node.Name)
	}
	return nil
}

func (m *MMD) getNode(name string, opts metav1.GetOptions) (*corev1.Node, error) {
	result := corev1.Node{}
	err := m.client.
		Get().
		Resource("nodes").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}
