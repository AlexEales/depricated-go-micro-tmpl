package k8s

import (
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client defines a interface to interact with kubernetes resources
type Client interface {
	corev1.CoreV1Interface
	WaitFor() PodWaiter
}

// NewClient creates a new k8s client using the kube config at the provided location
func NewClient(kubeconfig string) (Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &k8sClient{
		CoreV1Interface: clientset.CoreV1(),
	}, nil
}

// NewInClusterClient creates a new k8s client with a in-cluster configuration and context
func NewInClusterClient() (Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &k8sClient{
		CoreV1Interface: clientset.CoreV1(),
	}, nil
}

type k8sClient struct {
	corev1.CoreV1Interface
}

// WaitFor returns a new pod waiter to be used to wait for pod states and phases concurrently
func (c *k8sClient) WaitFor() PodWaiter {
	return NewPodWaiter(c)
}
