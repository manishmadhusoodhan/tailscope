package main

import (
	"fmt"
	"net/http"

	"tailscope/internal/api"
	"tailscope/internal/app"
	"tailscope/internal/kube"
)

func main() {

	fmt.Println("Starting TailScope")

	kubeClient, err := kube.New()
	if err != nil {
		panic(err)
	}

	application := app.New(kubeClient)

	server := api.NewServer(
		application,
		kubeClient,
	)

	fmt.Println("HTTP API listening on :8080")

	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		panic(err)
	}
}
