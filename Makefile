GO_DOCKER_RUN = docker run --rm -v $(shell pwd):/go/src/github.com/nginxinc/nginx-asg-sync -v $(shell pwd)/build_output:/build_output -w /go/src/github.com/nginxinc/nginx-asg-sync/cmd/sync
GOFLAGS ?= -mod=vendor

export DOCKER_BUILDKIT = 1

all: amazon centos7 ubuntu-xenial amazon2 ubuntu-bionic

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

clean:
	-rm -r build_output

.PHONY: test
