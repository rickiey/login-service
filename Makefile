.PHONY: clean
export CGO_ENABLED=0

SVC_NAME=login-service
my_d=$(shell pwd)
OUT_D = $(my_d)/bin

UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)

BUILD := $(shell date +%y%m%d.%H%M)
DATE := $(shell date +%Y%m%d)
#VERSION_NO :=$(shell git describe --always --tags `git rev-list --tags --max-count=1`)
BUILD_FLAGS = -ldflags="-w -X main.Version=$(VERSION_NO) main.GitHash=$(GITCOMMIT) main.Build=$(BUILD)"


# docker 私有镜像地址
REGISTRY_URL=mysqlserv.com/common
#TAG?=1.8.0.09_20190515
TAG := $(shell git describe --always --tags)
DOCKERFILE?=Dockerfile

# 没有找到version tag，就使用 BUILD
ifeq ($(VERSION_NO), $(GITCOMMIT))
    VERSION_NO = $BUILD
endif

DEV=dev

#默认本地是mac
GOOS = darwin
GOARCH = amd64
APP:=$(SVC_NAME)

#获取机器OS，ARCH
ifneq ($(UNAME_S),Darwin)
   GOOS = linux
endif

ifneq ($(UNAME_P),x86_64)
    GOARCH = 386
endif

#默认编译成本地
all: native

native: _build

_build:
	@echo "Build 1/3 : create bin dir"
	@mkdir -p $(OUT_D)

	@echo "Build 2/3 : build at  " $(my_d)/cmd/$(SVC_NAME)
	@cd $(my_d)/cmd/$(SVC_NAME)

	@echo "Build 3/3 : building..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o $(OUT_D)/$(APP) ./cmd/$(SVC_NAME)
	@echo " --->  done!\n"

#开发环境
dev: GOOS=linux
dev: GOARCH=amd64
dev: CGO_ENABLED=0
dev: _build

docker: dev
	@sleep 1s
	docker build --build-arg SVC_NAME=$(APP) -f $(DOCKERFILE) -t $(REGISTRY_URL)/$(APP):$(TAG) .

push:
	docker push $(REGISTRY_URL)/$(APP):$(TAG)

# k8s 启动（没用过，仅供参考）
up:
	kubectl -n develop set image deployments/$(SVC_NAME) $(SVC_NAME)=$(REGISTRY_URL)/$(APP):$(TAG)

clean:
	@rm -rf bin
