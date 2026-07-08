package kube

import (
	"bufio"
	"context"
	"io"

	v1 "k8s.io/api/core/v1"

	"tailscope/internal/model"
)

func (c *Client) streamPodLogs(
	ctx context.Context,
	namespace string,
	pod string,
	out chan<- model.RawLog,
) error {

	tail := int64(100)

	req := c.clientset.
		CoreV1().
		Pods(namespace).
		GetLogs(pod, &v1.PodLogOptions{
			Follow:    true,
			TailLines: &tail,
		})

	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}

	defer stream.Close()

	reader := bufio.NewReader(stream)

	for {

		select {

		case <-ctx.Done():
			return nil

		default:

		}

		line, err := reader.ReadString('\n')

		if len(line) > 0 {

			out <- model.RawLog{
				Namespace: namespace,
				Pod:       pod,
				Line:      line,
			}
		}

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
	}
}
