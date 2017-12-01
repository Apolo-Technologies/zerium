# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: zeriumd android ios zeriumd-cross swarm zvm all test clean
.PHONY: zeriumd-linux zeriumd-linux-386 zeriumd-linux-amd64 zeriumd-linux-mips64 zeriumd-linux-mips64le
.PHONY: zeriumd-linux-arm zeriumd-linux-arm-5 zeriumd-linux-arm-6 zeriumd-linux-arm-7 zeriumd-linux-arm64
.PHONY: zeriumd-darwin zeriumd-darwin-386 zeriumd-darwin-amd64
.PHONY: zeriumd-windows zeriumd-windows-386 zeriumd-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

zeriumd:
	build/env.sh go run build/ci.go install ./cmd/zeriumd
	@echo "Done building."
	@echo "Run \"$(GOBIN)/zeriumd\" to launch zeriumd."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/zeriumd.aar\" to use the library."

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

zeriumd-cross: zeriumd-linux zeriumd-darwin zeriumd-windows zeriumd-android zeriumd-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-*

zeriumd-linux: zeriumd-linux-386 zeriumd-linux-amd64 zeriumd-linux-arm zeriumd-linux-mips64 zeriumd-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-*

zeriumd-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/zeriumd
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep 386

zeriumd-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/zeriumd
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep amd64

zeriumd-linux-arm: zeriumd-linux-arm-5 zeriumd-linux-arm-6 zeriumd-linux-arm-7 zeriumd-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep arm

zeriumd-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/zeriumd
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep arm-5

zeriumd-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/zeriumd
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep arm-6

zeriumd-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/zeriumd
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep arm-7

zeriumd-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/zeriumd
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep arm64

zeriumd-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/zeriumd
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep mips

zeriumd-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/zeriumd
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep mipsle

zeriumd-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/zeriumd
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep mips64

zeriumd-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/zeriumd
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-linux-* | grep mips64le

zeriumd-darwin: zeriumd-darwin-386 zeriumd-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-darwin-*

zeriumd-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/zeriumd
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-darwin-* | grep 386

zeriumd-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/zeriumd
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-darwin-* | grep amd64

zeriumd-windows: zeriumd-windows-386 zeriumd-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-windows-*

zeriumd-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/zeriumd
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-windows-* | grep 386

zeriumd-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/zeriumd
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zeriumd-windows-* | grep amd64
