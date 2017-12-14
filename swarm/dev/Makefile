.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/apolo-technologies/zerium/cmd/swarm
	go build -o bin/zrmd     github.com/apolo-technologies/zerium/cmd/zrmd
	go build -o bin/bootnode github.com/apolo-technologies/zerium/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/apolo-technologies/zerium/swarm/...
