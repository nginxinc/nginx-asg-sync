[![Build Status](https://travis-ci.org/nginxinc/nginx-asg-sync.svg?branch=master)](https://travis-ci.org/nginxinc/nginx-asg-sync)  [![FOSSA Status](https://app.fossa.io/api/projects/custom%2B1062%2Fgithub.com%2Fnginxinc%2Fnginx-asg-sync.svg?type=shield)](https://app.fossa.io/projects/custom%2B1062%2Fgithub.com%2Fnginxinc%2Fnginx-asg-sync?ref=badge_shield)  [![Go Report Card](https://goreportcard.com/badge/github.com/nginxinc/nginx-asg-sync)](https://goreportcard.com/report/github.com/nginxinc/nginx-asg-sync)

# NGINX Plus Integration with Cloud Autoscaling -- nginx-asg-sync

**nginx-asg-sync** allows [NGINX Plus](https://www.nginx.com/products/) to discover instances (virtual machines) of a scaling group of a cloud provider. The following providers are supported:

* AWS [Auto Scaling groups](http://docs.aws.amazon.com/autoscaling/latest/userguide/WhatIsAutoScaling.html)
* Azure [Virtual Machine Scale Sets](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/)

When the number of instances changes, nginx-asg-sync adds the new instances to the NGINX Plus configuration and removes the terminated ones.

## How It Works
nginx-asg-sync is an agent process that runs on the same instance as NGINX Plus. It polls for changes to the backend instance groups via the Cloud Provider API.

When it sees that a scaling event has happened, it adds or removes the corresponding backend instances from the NGINX Plus configuration via the NGINX Plus API.

**Note:** nginx-asg-sync does not scale cloud scaling groups, it only gets the IP addresses of the instances in the groups.

In the example below (AWS), NGINX Plus is configured to load balance among the instances of two AWS Auto Scaling groups -- Backend One and Backend Two.
nginx-asg-sync, running on the same instance as NGINX Plus, ensures that whenever you scale the Auto Scaling groups, the corresponding instances are added (or removed) from the NGINX Plus configuration.

![nginx-asg-sync-architecture](https://cdn-1.wp.nginx.com/wp-content/uploads/2017/03/aws-auto-scaling-group-asg-sync.png)

Below you will find documentation on how to use nginx-asg-sync.

## Documentation
**Note:** the documentation for **the latest stable release** is available via a link in the description of the release. See the [releases page](https://github.com/nginxinc/nginx-asg-sync/releases).

**Contents:**
- [NGINX Plus Integration with Cloud Autoscaling -- nginx-asg-sync](#nginx-plus-integration-with-cloud-autoscaling----nginx-asg-sync)
  - [How It Works](#how-it-works)
  - [Documentation](#documentation)
  - [Supported Operating Systems](#supported-operating-systems)
  - [Installation](#installation)
    - [NGINX Plus Configuration](#nginx-plus-configuration)
    - [Configuration for Cloud Providers](#configuration-for-cloud-providers)
  - [Usage](#usage)
  - [Troubleshooting](#troubleshooting)
  - [Building a Software Package](#building-a-software-package)
  - [Support](#support)

## Supported Operating Systems

We provide packages for the following operating systems:

* Ubuntu: 14.04 (Trusty), 16.04 (Xenial), 18.04 (Bionic)
* CentOS/RHEL: 7
* Amazon Linux: 1, 2

Support for other operating systems can be added.

## Installation

1. Get a software package for your OS:
    * For a stable release, download a package from the [releases page](https://github.com/nginxinc/nginx-asg-sync/releases).
    * For the latest source code from the master branch, build a software package by following [these instructions](#building-a-software-package).
2. Install the package:
    * For Amazon Linux or CentOS/RHEL, run: `$ sudo rpm -i <package-name>.rpm`
    * For Ubuntu, run: `$ sudo dpkg -i <package-name>.deb`

### NGINX Plus Configuration

As an example, we configure NGINX Plus to load balance two groups of instances -- backend-group-one and backend-group-two. NGINX Plus routes requests to the appropriate group based on the request URI:

* Requests for /backend-one go to Backend One group.
* Requests for /backend-two go to Backend Two group.

This example corresponds to [the diagram](#how-it-works) at the top of this README.

```nginx
upstream backend-one {
   zone backend-one 64k;
   state /var/lib/nginx/state/backend-one.conf;
}

upstream backend-two {
   zone backend-two 64k;
   state /var/lib/nginx/state/backend-two.conf;
}

server {
   listen 80;

   status_zone backend;

   location /backend-one {
       proxy_set_header Host $host;
       proxy_pass http://backend-one;
   }

   location @hc-backend-one {
       internal;
       proxy_connect_timeout 1s;
       proxy_read_timeout 1s;
       proxy_send_timeout 1s;

       proxy_pass http://backend-one;
       health_check interval=1s mandatory;
   }

   location /backend-two {
       proxy_set_header Host $host;
       proxy_pass http://backend-two;
   }

   location @hc-backend-two {
       internal;
       proxy_connect_timeout 1s;
       proxy_read_timeout 1s;
       proxy_send_timeout 1s;

       proxy_pass http://backend-two;
       health_check interval=1s mandatory;
   }
}

server {
    listen 8080;

    location /api {
        api write=on;
    }

    location /dashboard.html {
      root /usr/share/nginx/html;
    }
}
```

* We declare two upstream groups – **backend-one** and **backend-two**, which correspond to our instance groups. However, we do not add any servers to the upstream groups, because the servers will be added by nginx-aws-sync. The `state` directive names the file where the dynamically configurable list of servers is stored, enabling it to persist across restarts of NGINX Plus.
* We define a virtual server that listens on port 80. NGINX Plus passes requests for **/backend-one** to the instances of the Backend One group, and requests for **/backend-two** to the instances of the Backend Two group.
* We define a second virtual server listening on port 8080 and configure the NGINX Plus API on it, which is required by nginx-asg-sync:
  * The API is available at **127.0.0.1:8080/api**

Because cloud provider APIs return the instances IP addresses before the instances are ready and/or provisioned, we recommend setting up mandatory active [healthchecks](http://nginx.org/en/docs/http/ngx_http_upstream_hc_module.html#health_check) for all upstream groups - **@hc-backend-one** and **@hc-backend-two**. This way, NGINX Plus won't pass any request to an instance that is still being provisioned or has been deleted recently.

Small timeouts ensure that a health check will fail fast if the backend instance is not healthy. Also, the mandatory parameter ensures NGINX Plus won't consider a newly added instance healthy until a health check passes.

### Configuration for Cloud Providers

See the example for your cloud provider: [AWS](examples/aws.md), [Azure](examples/azure.md).

## Usage

nginx-asg-sync runs as a system service and supports the start/stop/restart commands.

For Ubuntu 14.04 and Amazon Linux 1, run: `$ sudo start|stop|restart nginx-asg-sync`

For Ubuntu 16.04 and 18.04, CentOS7/RHEL7 and Amazon Linux 2, run: `$ sudo service nginx-asg-sync start|stop|restart`

## Troubleshooting

If nginx-asg-sync doesn’t work as expected, check its log file available at **/var/log/nginx-aws-sync/nginx-aws-sync.log**.

## Building a Software Package

You can compile nginx-asg-sync and build a software package using the provided Makefile. Before you start building a package, make sure that the following software is installed on your system:
* make
* Docker

To build a software package, run: `$ make <os>`
where `<os>` is the target OS. The following values are allowed:
* `amazon` for Amazon Linux 1
* `amazon2` for Amazon Linux 2
* `centos7` for CentOS7/RHEL7
* `ubuntu-trusty` for Ubuntu 14.04
* `ubuntu-xenial` for Ubuntu 16.04
* `ubuntu-bionic` for Ubuntu 18.04

If you run make without any arguments, it will build software packages for all supported OSes.

## Support

Support from the [NGINX Professional Services Team](https://www.nginx.com/services/) is available when using nginx-asg-sync.
