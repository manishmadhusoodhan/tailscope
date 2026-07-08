APP=tailscope
IMAGE=tailscope:latest

.PHONY: run build docker-build docker-run clean

run:
	go run ./cmd

build:
	go build -o bin/$(APP) ./cmd

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm \
		-p 8080:8080 \
		-v $(HOME)/.kube:/home/tailscope/.kube:ro \
		$(IMAGE)

clean:
	rm -rf bin