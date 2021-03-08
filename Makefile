# 主要配置项

# PROTO_PATHS = $(shell find pkg/ internal/ -name '*.proto' | xargs -I {} dirname {} | uniq)
PROTO_PATHS = $(shell find api/ proto/ -name '*.proto' | xargs -I {} dirname {} | uniq)


# CMD_TARGETS = $(notdir $(shell find cmd/* -maxdepth 0 -type d))


# target 实现

# .DEFAULT_GOAL := all
#
# .PHONY: deps all proto $(CMD_TARGETS) lint test test-e2e proto_path deps/corepb
#
# export PATH := $(shell pwd)/deps/:$(PATH)

# 依赖工具安装

# deps/protoc:
# 	bash scripts/get-protoc.sh
#
# deps/include/akali:
# 	bash scripts/get-akali-proto.sh
#
# deps/protoc-gen-go:
# 	export GOBIN=`pwd`/deps; cd; GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@v1.3.5
#
# deps/golangci-lint:
# 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b deps v1.33.0
#
# deps/external_services:
# 	export GOBIN=`pwd`/deps; cd; GO111MODULE=on go get -insecure git.code.oa.com/allocli/ci-tools/cmd/external_services@latest
#
# deps/ginkgo:
# 	export GOBIN=`pwd`/deps; cd; GO111MODULE=on go get github.com/onsi/ginkgo/ginkgo@latest
#
# deps/corepb:
# 	git submodule update --init --remote --merge
# 	mkdir -p deps/include/corepb 2>&1 >/dev/null
# 	cp `pwd`/submodules/game-proto/corepb/*.proto deps/include/corepb/
# 	go get git.woa.com/red/game-proto@master
#
# deps: deps/protoc deps/include/akali deps/corepb deps/protoc-gen-go deps/golangci-lint deps/external_services deps/ginkgo

# 构建应用

# all: deps $(CMD_TARGETS)
#
# $(CMD_TARGETS): proto
# 	CGO_ENABLED=0 go build -o bin/$@ ./cmd/$@

# 生成 protobuf 源码

proto: deps proto_path

proto_path: $(PROTO_PATHS)
	@$(foreach p,$^,protoc -I . --go_out=plugins=grpc,paths=source_relative:. $(wildcard $(p)/*.proto);)

# 测试

# 本地测试时按需在外部传递以下环境变量
# export SVC_DOMAIN ?= localhost
# export E2E_ID ?= localtest
# export E2E_GROUP ?= localtest
# export E2E_APP_VERSION ?= localtest

# lint: deps
# 	golangci-lint run ./...
#
# test: deps
# 	bash scripts/run-test.sh
#
# test-e2e: deps
# 	bash scripts/run-test-e2e.sh
