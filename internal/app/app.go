package app

import (
	"context"
	"sync"

	"tailscope/internal/kube"
	"tailscope/internal/output"
	"tailscope/internal/parser"
	"tailscope/internal/store"
)

type App struct {
	kube   *kube.Client
	parser *parser.Parser
	store  *store.Store

	mu     sync.Mutex
	cancel context.CancelFunc

	namespace string
	pods      []string
}

func New(
	kubeClient *kube.Client,
) *App {

	return &App{
		kube:   kubeClient,
		parser: parser.New(),
		store:  store.New(10000),
	}
}

func (a *App) Store() *store.Store {

	return a.store
}

// StartFollowing starts a new log stream for the selected namespace and pods.
// If another stream is already running, it is cancelled first.
func (a *App) StartFollowing(
	namespace string,
	pods []string,
) error {

	a.mu.Lock()
	defer a.mu.Unlock()

	// Stop existing stream
	if a.cancel != nil {
		a.cancel()
	}

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	a.cancel = cancel

	a.namespace = namespace
	a.pods = pods

	go func() {

		err := a.follow(
			ctx,
			namespace,
			pods,
		)

		if err != nil {
			output.PrintError(err)
		}

	}()

	return nil
}

// StopFollowing stops the current Kubernetes log stream.
func (a *App) StopFollowing() {

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cancel != nil {

		a.cancel()

		a.cancel = nil
	}
}

// StreamInfo returns the currently active stream configuration.
func (a *App) StreamInfo() StreamInfo {

	a.mu.Lock()
	defer a.mu.Unlock()

	return StreamInfo{
		Namespace: a.namespace,
		Pods:      a.pods,
		Running:   a.cancel != nil,
	}
}

func (a *App) follow(
	ctx context.Context,
	namespace string,
	pods []string,
) error {

	logs, err := a.kube.FollowPods(
		ctx,
		namespace,
		pods,
	)

	if err != nil {
		return err
	}

	for {

		select {

		case <-ctx.Done():

			return nil

		case raw, ok := <-logs:

			if !ok {
				return nil
			}

			entry := a.parser.Parse(raw)

			a.store.Add(entry)

			output.Print(entry)
		}
	}
}

type StreamInfo struct {
	Namespace string `json:"namespace"`

	Pods []string `json:"pods"`

	Running bool `json:"running"`
}
