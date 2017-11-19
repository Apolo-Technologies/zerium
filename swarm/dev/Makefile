.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/abt/zerium/cmd/swarm
	go build -o bin/gzrm     github.com/abt/zerium/cmd/gzrm
	go build -o bin/bootnode github.com/abt/zerium/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/abt/zerium/swarm/...
