package k8s

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client interface {
	corev1.CoreV1Interface
	WaitForPodStatus(string, string, v1.PodStatus) error
	WaitForPodStatuses(string, map[string]v1.PodStatus) error
}

func NewClient() Client {
	// FIXME: this client does not match the corev1 interface
	return &k8sClient{}
}

type k8sClient struct {
	*kubernetes.Clientset
}

func (c *k8sClient) WaitForPodStatus(namespace string, pod string, status v1.PodStatus) error {
	// TODO: https://bcreane.github.io/kubernetes/2018/05/10/golang-k8s-pod-ready.html
	return nil
}

func (c *k8sClient) WaitForPodStatuses(namespace string, podStatusMap map[string]v1.PodStatus) error {
	return nil
}
