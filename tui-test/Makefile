IMAGE ?= tui-test
TAG   ?= latest

.PHONY: build run docker-build docker-run

docker-build:
	docker build -t $(IMAGE):$(TAG) .

# -it allocates a TTY and keeps stdin open, which the TUI needs.
# TERM/COLORTERM are forwarded so colors match the host terminal.
docker-run:
	docker run --rm -it -e TERM=$${TERM:-xterm-256color} -e COLORTERM=$${COLORTERM} $(IMAGE):$(TAG)
