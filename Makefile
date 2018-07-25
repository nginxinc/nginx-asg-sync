GO_DOCKER_RUN = docker run --rm -v $(shell pwd):/go/src/github.com/nginxinc/nginx-asg-sync -w /go/src/github.com/nginxinc/nginx-asg-sync/cmd/sync
GOLANG_CONTAINER = golang:1.10

all: amazon centos7 ubuntu-trusty ubuntu-xenial

test:
	$(GO_DOCKER_RUN) $(GOLANG_CONTAINER) go test

compile: test
	$(GO_DOCKER_RUN) $(GOLANG_CONTAINER) go build -o /build/nginx-asg-sync

amazon: compile
	make -C os-packages/builders/amazon/
	docker run --rm -v $(shell pwd)/os-packages/rpm:/rpm -v $(shell pwd)/build:/build amazon-builder

centos7: compile
	make -C os-packages/builders/centos7/
	docker run --rm -v $(shell pwd)/os-packages/rpm:/rpm -v $(shell pwd)/build:/build centos7-builder

ubuntu-xenial: compile
	make -C os-packages/builders/ubuntu-xenial/
	docker run --rm -v $(shell pwd)/os-packages/debian:/debian -v $(shell pwd)/build:/build ubuntu-xenial-builder

ubuntu-trusty: compile
	make -C os-packages/builders/ubuntu-trusty/
	docker run --rm -v $(shell pwd)/os-packages/debian:/debian -v $(shell pwd)/build:/build ubuntu-trusty-builder

clean:
	-rm -r build

.PHONY: test
