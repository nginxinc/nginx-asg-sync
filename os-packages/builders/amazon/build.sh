#!/bin/bash

mkdir -p ~/rpmbuild
cp -r rpm/* ~/rpmbuild/
rpmbuild -bb ~/rpmbuild/SPECS/nginx-asg-sync.spec
cp ~/rpmbuild/RPMS/x86_64/*.rpm /build
