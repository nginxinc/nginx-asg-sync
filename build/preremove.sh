#!/bin/sh

systemctl --no-reload disable nginx-asg-sync.service >/dev/null 2>&1 || :
systemctl stop nginx-asg-sync.service >/dev/null 2>&1 || :
