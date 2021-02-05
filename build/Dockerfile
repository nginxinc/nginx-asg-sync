ARG CONTAINER_VERSION
ARG GOLANG_CONTAINER=golang:1.15

FROM $GOLANG_CONTAINER AS builder
WORKDIR /go/src/github.com/nginxinc/nginx-asg-sync/cmd/sync
COPY . /go/src/github.com/nginxinc/nginx-asg-sync/

RUN go build -o /nginx-asg-sync

#---------------------------------------------------------------------------------------------

FROM ${CONTAINER_VERSION} as rpm_based

RUN yum install -y rpmdevtools
ADD build/package/builders/rpm_based/build.sh /

COPY --from=builder /nginx-asg-sync /nginx-asg-sync

ENTRYPOINT ["/build.sh"]

#---------------------------------------------------------------------------------------------

FROM ${CONTAINER_VERSION} as deb_based

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install debhelper dh-systemd -y
ADD build/package/builders/deb_based/build.sh /

COPY --from=builder /nginx-asg-sync /nginx-asg-sync

ARG PACKAGE_VERSION
ARG OS_VERSION
ENV PACKAGE_VERSION=${PACKAGE_VERSION}
ENV OS_VERSION=${OS_VERSION}

ENTRYPOINT ["/build.sh"]
