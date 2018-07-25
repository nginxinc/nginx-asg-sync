FROM amazonlinux:latest

RUN yum install -y rpmdevtools
ADD build.sh /

ENTRYPOINT ["/build.sh"]
