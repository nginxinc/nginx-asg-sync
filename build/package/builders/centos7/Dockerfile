FROM centos:7

RUN yum install -y rpmdevtools
ADD build.sh /

ENTRYPOINT ["/build.sh"]
