#!/bin/bash

mkdir -p ~/nginx-asg-sync-0.1/
cp -r /debian ~/nginx-asg-sync-0.1/
cd ~/nginx-asg-sync-0.1/
sed -i 's/%%CODENAME%%/trusty/g' debian/changelog
rm debian/nginx-asg-sync.service
dpkg-buildpackage -us -uc
cd ..
mv *.deb /build/
