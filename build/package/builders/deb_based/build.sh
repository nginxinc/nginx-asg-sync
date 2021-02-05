#!/bin/bash

# package_name=nginx-asg-sync-${PACKAGE_VERSION}

export SYSTEMD='--with=systemd'
mkdir -p ~/nginx-asg-sync/
cp -r /debian ~/nginx-asg-sync/
cd ~/nginx-asg-sync/
sed -i "s/%%CODENAME%%/${OS_VERSION}/g" debian/changelog
rm debian/nginx-asg-sync.upstart
dpkg-buildpackage -us -uc
cd ..
mv *.deb /build_output/
