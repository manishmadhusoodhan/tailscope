package kube

import (
	"context"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func New() (*Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		kubeconfig = filepath.Join(
			homedir.HomeDir(),
			".kube",
			"config",
		)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func (c *Client) Namespaces(
	ctx context.Context,
) ([]string, error) {

	list, err :=
		c.clientset.CoreV1().
			Namespaces().
			List(
				ctx,
				metav1.ListOptions{},
			)

	if err != nil {

		return nil, err

	}

	result :=
		make([]string, 0, len(list.Items))

	for _, ns := range list.Items {

		result =
			append(
				result,
				ns.Name,
			)

	}

	return result, nil

}

func (c *Client) PodsForNamespace(
	ctx context.Context,
	namespace string,
) ([]string, error) {

	list, err := c.clientset.
		CoreV1().
		Pods(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	pods := make([]string, 0, len(list.Items))

	for _, pod := range list.Items {
		pods = append(pods, pod.Name)
	}

	return pods, nil
}
