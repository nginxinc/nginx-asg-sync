#!/bin/bash

export SYSTEMD='--with=systemd'
mkdir -p ~/nginx-asg-sync-0.2/
cp -r /debian ~/nginx-asg-sync-0.2/
cd ~/nginx-asg-sync-0.2/
sed -i 's/%%CODENAME%%/bionic/g' debian/changelog
rm debian/nginx-asg-sync.upstart
dpkg-buildpackage -us -uc
cd ..
mv *.deb /build_output/
