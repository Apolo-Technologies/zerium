# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gabt android ios gabt-cross swarm evm all test clean
.PHONY: gabt-linux gabt-linux-386 gabt-linux-amd64 gabt-linux-mips64 gabt-linux-mips64le
.PHONY: gabt-linux-arm gabt-linux-arm-5 gabt-linux-arm-6 gabt-linux-arm-7 gabt-linux-arm64
.PHONY: gabt-darwin gabt-darwin-386 gabt-darwin-amd64
.PHONY: gabt-windows gabt-windows-386 gabt-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gabt:
	build/env.sh go run build/ci.go install ./cmd/gabt
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gabt\" to launch gabt."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gabt.aar\" to use the library."

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

gabt-cross: gabt-linux gabt-darwin gabt-windows gabt-android gabt-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gabt-*

gabt-linux: gabt-linux-386 gabt-linux-amd64 gabt-linux-arm gabt-linux-mips64 gabt-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-*

gabt-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gabt
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep 386

gabt-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gabt
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep amd64

gabt-linux-arm: gabt-linux-arm-5 gabt-linux-arm-6 gabt-linux-arm-7 gabt-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep arm

gabt-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gabt
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep arm-5

gabt-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gabt
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep arm-6

gabt-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gabt
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep arm-7

gabt-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gabt
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep arm64

gabt-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gabt
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep mips

gabt-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gabt
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep mipsle

gabt-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gabt
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep mips64

gabt-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gabt
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gabt-linux-* | grep mips64le

gabt-darwin: gabt-darwin-386 gabt-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gabt-darwin-*

gabt-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gabt
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-darwin-* | grep 386

gabt-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gabt
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-darwin-* | grep amd64

gabt-windows: gabt-windows-386 gabt-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gabt-windows-*

gabt-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gabt
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-windows-* | grep 386

gabt-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gabt
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gabt-windows-* | grep amd64
