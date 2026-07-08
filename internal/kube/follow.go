package kube

import (
	"context"
	"sync"

	"tailscope/internal/model"
)

func (c *Client) FollowPods(
	ctx context.Context,
	namespace string,
	pods []string,
) (<-chan model.RawLog, error) {

	out := make(chan model.RawLog, 1000)

	var wg sync.WaitGroup

	for _, pod := range pods {

		wg.Add(1)

		go func(name string) {
			defer wg.Done()

			_ = c.streamPodLogs(
				ctx,
				namespace,
				name,
				out,
			)
		}(pod)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out, nil
}
