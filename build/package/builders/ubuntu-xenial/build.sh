#!/bin/bash

export SYSTEMD='--with=systemd'
mkdir -p ~/nginx-asg-sync-0.1/
cp -r /debian ~/nginx-asg-sync-0.1/
cd ~/nginx-asg-sync-0.1/
sed -i 's/%%CODENAME%%/xenial/g' debian/changelog
rm debian/nginx-asg-sync.upstart
dpkg-buildpackage -us -uc
cd ..
mv *.deb /build_output/
