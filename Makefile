GO_DOCKER_RUN = docker run --rm -v $(shell pwd):/go/src/github.com/nginxinc/nginx-asg-sync -v $(shell pwd)/build_output:/build_output -w /go/src/github.com/nginxinc/nginx-asg-sync/cmd/sync
GOLANG_CONTAINER = golang:1.12
BUILD_IN_CONTAINER = 1

all: amazon centos7 ubuntu-xenial amazon2 ubuntu-bionic

test:
ifeq ($(BUILD_IN_CONTAINER),1)
	$(GO_DOCKER_RUN) $(GOLANG_CONTAINER) go test
else
	go test ./...
endif

lint:
	golangci-lint run

compile: test
ifeq ($(BUILD_IN_CONTAINER),1)
	$(GO_DOCKER_RUN) $(GOLANG_CONTAINER) go build -o /build_output/nginx-asg-sync
else
	go build -o ./build_output/nginx-asg-sync github.com/nginxinc/nginx-asg-sync/cmd/sync
endif

amazon: compile
	make -C build/package/builders/amazon/
	docker run --rm -v $(shell pwd)/build/package/rpm:/rpm -v $(shell pwd)/build_output:/build_output amazon-builder

amazon2: compile
	make -C build/package/builders/amazon2/
	docker run --rm -v $(shell pwd)/build/package/rpm:/rpm -v $(shell pwd)/build_output:/build_output amazon2-builder

centos7: compile
	make -C build/package/builders/centos7/
	docker run --rm -v $(shell pwd)/build/package/rpm:/rpm -v $(shell pwd)/build_output:/build_output centos7-builder

ubuntu-xenial: compile
	make -C build/package/builders/ubuntu-xenial/
	docker run --rm -v $(shell pwd)/build/package/debian:/debian -v $(shell pwd)/build_output:/build_output ubuntu-xenial-builder

ubuntu-bionic: compile
	make -C build/package/builders/ubuntu-bionic/
	docker run --rm -v $(shell pwd)/build/package/debian:/debian -v $(shell pwd)/build_output:/build_output ubuntu-bionic-builder

clean:
	-rm -r build_output

.PHONY: test
