package kube

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (c *Client) Deployments(namespace string) ([]appsv1.Deployment, error) {

	list, err := c.clientset.
		AppsV1().
		Deployments(namespace).
		List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (c *Client) Deployment(namespace, name string) (*appsv1.Deployment, error) {

	return c.clientset.
		AppsV1().
		Deployments(namespace).
		Get(context.Background(), name, metav1.GetOptions{})
}

func (c *Client) PodsForDeployment(
	namespace string,
	deployment string,
) ([]v1.Pod, error) {

	dep, err := c.Deployment(namespace, deployment)
	if err != nil {
		return nil, err
	}

	selector := labels.Set(
		dep.Spec.Selector.MatchLabels,
	).String()

	list, err := c.clientset.
		CoreV1().
		Pods(namespace).
		List(context.Background(), metav1.ListOptions{
			LabelSelector: selector,
		})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (c *Client) ListDeployments(
	ctx context.Context,
	namespace string,
) ([]appsv1.Deployment, error) {

	deployments, err :=
		c.clientset.
			AppsV1().
			Deployments(namespace).
			List(
				ctx,
				metav1.ListOptions{},
			)

	if err != nil {
		return nil, err
	}

	return deployments.Items, nil
}
