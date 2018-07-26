FROM amazonlinux:2018.03.0.20180424

RUN yum install -y rpmdevtools
ADD build.sh /

ENTRYPOINT ["/build.sh"]
