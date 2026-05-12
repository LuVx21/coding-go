.PHONY: ai kv

ai:
	go build -ldflags="-w -s" -v -o $$(readlink -f $$(which ai)) ./ai

docker-tags:
	go build -ldflags="-w -s" -v -o $$(readlink -f $$(which docker-tags)) ./docker

kv:
	go build -ldflags="-w -s" -v -o $$(readlink -f $$(which kv)) ./kv
