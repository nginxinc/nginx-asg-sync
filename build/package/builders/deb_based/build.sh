#!/bin/bash

mkdir -p ~/nginx-asg-sync/
cp -r /debian ~/nginx-asg-sync/
cd ~/nginx-asg-sync/
dpkg-buildpackage -us -uc
cd ..
mv *.deb /build_output/
