FROM amazonlinux:2

RUN yum install -y rpmdevtools
ADD build.sh /

ENTRYPOINT ["/build.sh"]
