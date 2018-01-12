# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: zaed android ios zaed-cross swarm evm all test clean
.PHONY: zaed-linux zaed-linux-386 zaed-linux-amd64 zaed-linux-mips64 zaed-linux-mips64le
.PHONY: zaed-linux-arm zaed-linux-arm-5 zaed-linux-arm-6 zaed-linux-arm-7 zaed-linux-arm64
.PHONY: zaed-darwin zaed-darwin-386 zaed-darwin-amd64
.PHONY: zaed-windows zaed-windows-386 zaed-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

zaed:
	build/env.sh go run build/ci.go install ./cmd/zaed
	@echo "Done building."
	@echo "Run \"$(GOBIN)/zaed\" to launch zaed."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/zaed.aar\" to use the library."

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

zaed-cross: zaed-linux zaed-darwin zaed-windows zaed-android zaed-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/zaed-*

zaed-linux: zaed-linux-386 zaed-linux-amd64 zaed-linux-arm zaed-linux-mips64 zaed-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-*

zaed-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/zaed
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep 386

zaed-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/zaed
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep amd64

zaed-linux-arm: zaed-linux-arm-5 zaed-linux-arm-6 zaed-linux-arm-7 zaed-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep arm

zaed-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/zaed
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep arm-5

zaed-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/zaed
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep arm-6

zaed-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/zaed
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep arm-7

zaed-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/zaed
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep arm64

zaed-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/zaed
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep mips

zaed-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/zaed
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep mipsle

zaed-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/zaed
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep mips64

zaed-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/zaed
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/zaed-linux-* | grep mips64le

zaed-darwin: zaed-darwin-386 zaed-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/zaed-darwin-*

zaed-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/zaed
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-darwin-* | grep 386

zaed-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/zaed
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-darwin-* | grep amd64

zaed-windows: zaed-windows-386 zaed-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/zaed-windows-*

zaed-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/zaed
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-windows-* | grep 386

zaed-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/zaed
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zaed-windows-* | grep amd64
