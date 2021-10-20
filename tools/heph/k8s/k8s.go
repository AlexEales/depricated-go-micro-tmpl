package k8s

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client defines a interface to interact with kubernetes resources
type Client interface {
	corev1.CoreV1Interface
	WaitForPodToBeReady(context.Context, string, string, time.Duration) error
	WaitForPodsToBeReady(context.Context, string, []string, time.Duration) error
}

// TODO: The creation of the clientset should probably be done in here
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

// WaitForPodToBeReady waits the specified amount of time, polling every second, for the pod in the specified namespace
// to be ready
func (c *k8sClient) WaitForPodToBeReady(ctx context.Context, namespace, pod string, timeout time.Duration) error {
	return wait.PollImmediateWithContext(ctx, time.Second, timeout, func(ctx context.Context) (done bool, err error) {
		pod, err := c.Pods(namespace).Get(ctx, pod, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		conditions := pod.Status.Conditions
		for _, condition := range conditions {
			if condition.Type == v1.PodReady && condition.Status == v1.ConditionTrue {
				return true, nil
			}
		}

		return false, nil
	})
}

// WaitForPodsToBeReady waits the specified amount of time, polling every second, for the pods in the specified namespace
// to be ready
func (c *k8sClient) WaitForPodsToBeReady(ctx context.Context, namespace string, pods []string, timeout time.Duration) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, pod := range pods {
		pod := pod

		g.Go(func() error {
			return c.WaitForPodToBeReady(ctx, namespace, pod, timeout)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
