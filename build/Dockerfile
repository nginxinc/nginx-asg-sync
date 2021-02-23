ARG CONTAINER_VERSION=ubuntu:groovy
ARG OS_TYPE=deb_based
ARG PACKAGE_VERSION=0.4-1

FROM golang:1.16-alpine AS builder
WORKDIR /go/src/github.com/nginxinc/nginx-asg-sync/cmd/sync
COPY . /go/src/github.com/nginxinc/nginx-asg-sync/

RUN CGO_ENABLED=0 GO111MODULE=on GOFLAGS='-mod=vendor' go build -installsuffix cgo -o /nginx-asg-sync

#---------------------------------------------------------------------------------------------

FROM ${CONTAINER_VERSION} as rpm_based

RUN yum install -y rpmdevtools
ADD build/package/builders/rpm_based/build.sh /

ENTRYPOINT ["/build.sh"]

#---------------------------------------------------------------------------------------------

FROM ${CONTAINER_VERSION} as deb_based

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install debhelper dh-systemd -y
ADD build/package/builders/deb_based/build.sh /

ARG PACKAGE_VERSION
ARG CONTAINER_VERSION
ENV PACKAGE_VERSION=${PACKAGE_VERSION}
ENV CONTAINER_VERSION=${CONTAINER_VERSION}

ENTRYPOINT ["/build.sh"]

#---------------------------------------------------------------------------------------------

FROM ${OS_TYPE} as container

COPY --from=builder /nginx-asg-sync /nginx-asg-sync


#---------------------------------------------------------------------------------------------

FROM ${OS_TYPE} as local

COPY nginx-asg-sync /nginx-asg-sync
