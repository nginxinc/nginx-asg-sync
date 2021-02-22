#!/bin/bash

package_name=nginx-asg-sync-${PACKAGE_VERSION}

export SYSTEMD='--with=systemd'
export OS_VERSION=${CONTAINER_VERSION#*:}

mkdir -p ~/${package_name}/
cp -r /debian ~/${package_name}/
cd ~/${package_name}/
sed -i "s/%%CODENAME%%/${OS_VERSION}/g" debian/changelog
rm debian/nginx-asg-sync.upstart
dpkg-buildpackage -us -uc
cd ..
mv *.deb /build_output/
