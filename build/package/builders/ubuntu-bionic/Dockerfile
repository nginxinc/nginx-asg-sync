FROM ubuntu:bionic

RUN apt-get update && apt-get install debhelper dh-systemd -y
ADD build.sh /

ENTRYPOINT ["/build.sh"]
