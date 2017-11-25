# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gzrm android ios gzrm-cross swarm evm all test clean
.PHONY: gzrm-linux gzrm-linux-386 gzrm-linux-amd64 gzrm-linux-mips64 gzrm-linux-mips64le
.PHONY: gzrm-linux-arm gzrm-linux-arm-5 gzrm-linux-arm-6 gzrm-linux-arm-7 gzrm-linux-arm64
.PHONY: gzrm-darwin gzrm-darwin-386 gzrm-darwin-amd64
.PHONY: gzrm-windows gzrm-windows-386 gzrm-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gzrm:
	build/env.sh go run build/ci.go install ./cmd/gzrm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gzrm\" to launch gzrm."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gzrm.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gzrm.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/jteeuwen/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go install ./cmd/abigen

# Cross Compilation Targets (xgo)

gzrm-cross: gzrm-linux gzrm-darwin gzrm-windows gzrm-android gzrm-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-*

gzrm-linux: gzrm-linux-386 gzrm-linux-amd64 gzrm-linux-arm gzrm-linux-mips64 gzrm-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-*

gzrm-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gzrm
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep 386

gzrm-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gzrm
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep amd64

gzrm-linux-arm: gzrm-linux-arm-5 gzrm-linux-arm-6 gzrm-linux-arm-7 gzrm-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep arm

gzrm-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gzrm
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep arm-5

gzrm-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gzrm
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep arm-6

gzrm-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gzrm
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep arm-7

gzrm-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gzrm
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep arm64

gzrm-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gzrm
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep mips

gzrm-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gzrm
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep mipsle

gzrm-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gzrm
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep mips64

gzrm-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gzrm
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-linux-* | grep mips64le

gzrm-darwin: gzrm-darwin-386 gzrm-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-darwin-*

gzrm-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gzrm
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-darwin-* | grep 386

gzrm-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gzrm
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-darwin-* | grep amd64

gzrm-windows: gzrm-windows-386 gzrm-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-windows-*

gzrm-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gzrm
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-windows-* | grep 386

gzrm-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gzrm
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gzrm-windows-* | grep amd64
