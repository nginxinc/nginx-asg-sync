FROM ubuntu:xenial

RUN apt-get update && apt-get install debhelper dh-systemd -y
ADD build.sh /

ENTRYPOINT ["/build.sh"]
