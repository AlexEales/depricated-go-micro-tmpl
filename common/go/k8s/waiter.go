package k8s

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var log = logrus.New()

// PodWaiter defines the interface for a pod listener which waits for pods to reach the
// specified state or phase.
type PodWaiter interface {
	PodCondition(string, string, v1.PodConditionType) PodWaiter
	PodsToFulfillCondition(string, []string, v1.PodConditionType) PodWaiter
	PodPhase(string, string, v1.PodPhase) PodWaiter
	PodsToReachPhase(string, []string, v1.PodPhase) PodWaiter
	Wait(context.Context) error
	WithPollInterval(time.Duration) PodWaiter
	WithTimeout(time.Duration) PodWaiter
}

// NewPodWaiter returns a new pod waiter with a default poll interval of 5 seconds and
// a timeout of 2 minutes.
func NewPodWaiter(client corev1.CoreV1Interface) PodWaiter {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return &waiter{
		CoreV1Interface: client,
		conditionFns:    []wait.ConditionWithContextFunc{},
		pollInterval:    5 * time.Second,
		timeout:         2 * time.Minute,
		waitMap:         make(map[string]bool),
	}
}

type waiter struct {
	corev1.CoreV1Interface
	conditionFns []wait.ConditionWithContextFunc
	pollInterval time.Duration
	timeout      time.Duration
	waitMap      map[string]bool
}

// PodCondition adds a wait for a specified pod in the namespace to reach a specified condition (status).
func (w *waiter) PodCondition(namespace, pod string, condition v1.PodConditionType) PodWaiter {
	w.conditionFns = append(w.conditionFns, w.podConditionReachedFn(namespace, pod, condition))
	return w
}

// PodsToFulfillCondition adds a wait for a specified pods in the namespace to reach a specified condition (status).
func (w *waiter) PodsToFulfillCondition(namespace string, pods []string, condition v1.PodConditionType) PodWaiter {
	for _, pod := range pods {
		w.conditionFns = append(w.conditionFns, w.podConditionReachedFn(namespace, pod, condition))
	}

	return w
}

func (w *waiter) podConditionReachedFn(namespace, pod string, condition v1.PodConditionType) wait.ConditionWithContextFunc {
	resourceName := fmt.Sprintf("%s/%s", namespace, pod)
	w.waitMap[resourceName] = true

	return func(ctx context.Context) (done bool, err error) {
		pod, err := w.Pods(namespace).Get(ctx, pod, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		conditions := pod.Status.Conditions
		for _, c := range conditions {
			if c.Type == condition && c.Status == v1.ConditionTrue {
				delete(w.waitMap, resourceName)
				return true, nil
			}
		}

		return false, nil
	}
}

// PodPhase adds a wait for a specified pod in the namespace to reach a specified phase.
func (w *waiter) PodPhase(namespace, pod string, phase v1.PodPhase) PodWaiter {
	w.conditionFns = append(w.conditionFns, w.podPhaseReachedFn(namespace, pod, phase))
	return w
}

// PodsToReachPhase adds a wait for a specified pods in the namespace to reach a specified phase.
func (w *waiter) PodsToReachPhase(namespace string, pods []string, phase v1.PodPhase) PodWaiter {
	for _, pod := range pods {
		w.conditionFns = append(w.conditionFns, w.podPhaseReachedFn(namespace, pod, phase))
	}

	return w
}

func (w *waiter) podPhaseReachedFn(namespace, pod string, phase v1.PodPhase) wait.ConditionWithContextFunc {
	resourceName := fmt.Sprintf("%s/%s", namespace, pod)
	w.waitMap[resourceName] = true

	return func(ctx context.Context) (done bool, err error) {
		pod, err := w.Pods(namespace).Get(ctx, pod, metav1.GetOptions{})
		if err != nil {
			// If the pod isn't present ignore the error as it might still be creating
			if statusErr, ok := err.(*k8serrors.StatusError); ok {
				if statusErr.ErrStatus.Code == http.StatusNotFound {
					return false, nil
				}
			}
			return false, err
		}

		if pod.Status.Phase == phase {
			delete(w.waitMap, resourceName)
			return true, nil
		}

		return false, nil
	}
}

// Wait starts waiting on the conditions defined earlier in the chain, printing the pods
// still being waited on at the specified intervals.
func (w *waiter) Wait(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, conditionFn := range w.conditionFns {
		conditionFn := conditionFn

		g.Go(func() error {
			return wait.PollImmediateWithContext(ctx, w.pollInterval, w.timeout, conditionFn)
		})
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go w.printPodsBeingWaitedOn(cancelCtx)

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *waiter) printPodsBeingWaitedOn(ctx context.Context) {
	for range time.NewTicker(w.pollInterval).C {
		select {
		case <-ctx.Done():
			return
		default:
			waitingOn := w.getPodsBeingWaitedOn()
			log.Infof(
				"waiting on %d pod(s) to reach the required status/phase: {%s}",
				len(waitingOn),
				strings.Join(waitingOn, ", "),
			)
		}
	}
}

func (w *waiter) getPodsBeingWaitedOn() []string {
	pods := make([]string, len(w.waitMap))

	i := 0
	for pod, _ := range w.waitMap {
		pods[i] = pod
		i++
	}

	return pods
}

// WithPollInterval sets the waiters poll interval to the specified value.
func (w *waiter) WithPollInterval(interval time.Duration) PodWaiter {
	w.pollInterval = interval
	return w
}

// WithTimeout sets the waiters timeout to the specified value.
func (w *waiter) WithTimeout(timeout time.Duration) PodWaiter {
	w.timeout = timeout
	return w
}
