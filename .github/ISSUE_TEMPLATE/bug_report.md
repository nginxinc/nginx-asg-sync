---
name: Bug report
about: Create a report to help us improve

---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**

Provide the following files as part of the bug report

- NGINX Plus configuration. Run `sudo nginx -T` to gather the full configuration
- nginx-asg-sync configuration from `/etc/nginx/config.yaml`

Steps to reproduce the behavior, such as:

1. Scale from 2 to 5 EC2 instances
2. New instances not added to nginx.conf
3. See error in `/var/log/nginx-asg-sync/nginx-asg-sync.log`

**Expected behavior**
A clear and concise description of what you expected to happen.

**Your environment**

- Version of nginx-asg-sync
- Version of NGINX Plus
- Version of the OS

**Additional context**
Add any other context about the problem here. Any log files you want to share.
