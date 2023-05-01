#!/bin/sh

systemctl preset nginx-asg-sync.service >/dev/null 2>&1 || :
systemctl enable nginx-asg-sync.service >/dev/null 2>&1 || :
