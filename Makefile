GO_DOCKER_RUN = docker run --rm -v $(shell pwd):/go/src/github.com/nginxinc/nginx-asg-sync -v $(shell pwd)/build_output:/build_output -w /go/src/github.com/nginxinc/nginx-asg-sync/cmd/sync
GOFLAGS ?= -mod=vendor

export DOCKER_BUILDKIT = 1

all: amazon centos7 ubuntu-xenial amazon2 ubuntu-bionic ubuntu-focal ubuntu-groovy

.PHONY: test
test:
	GO111MODULE=on GOFLAGS='$(GOFLAGS)' go test ./...

lint:
	golangci-lint run

amazon:
	docker build -t amazon-builder --target rpm_based --build-arg CONTAINER_VERSION=amazonlinux:1 -f build/Dockerfile .
	docker run --rm -v $(shell pwd)/build/package/rpm:/rpm -v $(shell pwd)/build_output:/build_output amazon-builder

amazon2:
	docker build -t amazon2-builder --target rpm_based --build-arg CONTAINER_VERSION=amazonlinux:2 -f build/Dockerfile .
	docker run --rm -v $(shell pwd)/build/package/rpm:/rpm -v $(shell pwd)/build_output:/build_output amazon2-builder

centos7:
	docker build -t centos7-builder --target rpm_based --build-arg CONTAINER_VERSION=centos:7 -f build/Dockerfile .
	docker run --rm -v $(shell pwd)/build/package/rpm:/rpm -v $(shell pwd)/build_output:/build_output centos7-builder

ubuntu-xenial:
	docker build -t ubuntu-xenial-builder --target deb_based --build-arg CONTAINER_VERSION=ubuntu:xenial --build-arg OS_VERSION=xenial -f build/Dockerfile .
	docker run --rm -v $(shell pwd)/build/package/debian:/debian -v $(shell pwd)/build_output:/build_output ubuntu-xenial-builder

ubuntu-bionic:
	docker build -t ubuntu-bionic-builder --target deb_based --build-arg CONTAINER_VERSION=ubuntu:bionic --build-arg OS_VERSION=bionic -f build/Dockerfile .
	docker run --rm -v $(shell pwd)/build/package/debian:/debian -v $(shell pwd)/build_output:/build_output ubuntu-bionic-builder

ubuntu-focal:
	docker build -t ubuntu-focal-builder --target deb_based --build-arg CONTAINER_VERSION=ubuntu:focal --build-arg OS_VERSION=focal -f build/Dockerfile .
	docker run --rm -v $(shell pwd)/build/package/debian:/debian -v $(shell pwd)/build_output:/build_output ubuntu-focal-builder

ubuntu-groovy:
	docker build -t ubuntu-groovy-builder --target deb_based --build-arg CONTAINER_VERSION=ubuntu:groovy --build-arg OS_VERSION=groovy -f build/Dockerfile .
	docker run --rm -v $(shell pwd)/build/package/debian:/debian -v $(shell pwd)/build_output:/build_output ubuntu-groovy-builder

.PHONY: clean
clean:
	-rm -r build_output

.PHONY: deps
deps:
	@go mod tidy && go mod verify && go mod vendor

.PHONY: clean-cache
clean-cache:
	@go clean -modcache
