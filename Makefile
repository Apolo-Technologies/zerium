# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: zrmd android ios zrmd-cross swarm zvm all test clean
.PHONY: zrmd-linux zrmd-linux-386 zrmd-linux-amd64 zrmd-linux-mips64 zrmd-linux-mips64le
.PHONY: zrmd-linux-arm zrmd-linux-arm-5 zrmd-linux-arm-6 zrmd-linux-arm-7 zrmd-linux-arm64
.PHONY: zrmd-darwin zrmd-darwin-386 zrmd-darwin-amd64
.PHONY: zrmd-windows zrmd-windows-386 zrmd-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

zrmd:
	build/env.sh go run build/ci.go install ./cmd/zrmd
	@echo "Done building."
	@echo "Run \"$(GOBIN)/zrmd\" to launch zrmd."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/zrmd.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

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

zrmd-cross: zrmd-linux zrmd-darwin zrmd-windows zrmd-android zrmd-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-*

zrmd-linux: zrmd-linux-386 zrmd-linux-amd64 zrmd-linux-arm zrmd-linux-mips64 zrmd-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-*

zrmd-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/zrmd
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep 386

zrmd-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/zrmd
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep amd64

zrmd-linux-arm: zrmd-linux-arm-5 zrmd-linux-arm-6 zrmd-linux-arm-7 zrmd-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep arm

zrmd-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/zrmd
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep arm-5

zrmd-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/zrmd
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep arm-6

zrmd-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/zrmd
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep arm-7

zrmd-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/zrmd
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep arm64

zrmd-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/zrmd
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep mips

zrmd-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/zrmd
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep mipsle

zrmd-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/zrmd
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep mips64

zrmd-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/zrmd
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-linux-* | grep mips64le

zrmd-darwin: zrmd-darwin-386 zrmd-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-darwin-*

zrmd-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/zrmd
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-darwin-* | grep 386

zrmd-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/zrmd
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-darwin-* | grep amd64

zrmd-windows: zrmd-windows-386 zrmd-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-windows-*

zrmd-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/zrmd
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-windows-* | grep 386

zrmd-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/zrmd
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/zrmd-windows-* | grep amd64
