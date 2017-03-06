FROM ubuntu:trusty

RUN apt-get update && apt-get install debhelper -y
ADD build.sh /

ENTRYPOINT ["/build.sh"]
