# NGINX Plus Integration with AWS Auto Scaling groups -- nginx-asg-sync

**nginx-asg-sync** allows [NGINX Plus](https://www.nginx.com/products/) to support scaling when load balancing [AWS Auto Scaling groups](http://docs.aws.amazon.com/autoscaling/latest/userguide/WhatIsAutoScaling.html): when the number of instances in an Auto Scaling group changes, nginx-asg-sync adds the new instances to the NGINX Plus configuration and removes the terminated ones.

More details on this solution are available in the blog post [Load Balancing AWS Auto Scaling Groups with NGINX Plus](https://www.nginx.com/blog/load-balancing-aws-auto-scaling-groups-nginx-plus/).

Below you will find instructions on how to use nginx-asg-sync.

## Contents

1. [Supported Operating Systems](#supported-operating-systems)
1. [Setting up Access to the AWS API](#setting-up-access-to-the-aws-api)
1. [Installation](#installation)
1. [Configuration](#configuration)
1. [Usage](#usage)
1. [Troubleshooting](#troubleshooting)
1. [Building a Software Package](#building-a-software-package)
1. [Support](#support)

## Supported Operating Systems

We provide packages for the following operating systems:

* Ubuntu: 14.04 (Trusty), 16.04 (Xenial)
* CentOS/RHEL: 7
* Amazon Linux

Support for other operating systems can be added.

## Setting up Access to the AWS API

nginx-asg-sync uses the AWS API to get the list of IP addresses of the instances of an Auto Scaling group. To access the AWS API, nginx-asg-sync must have credentials. To provide credentials to nginx-asg-sync:

1. [Create an IAM role](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) and attach the predefined `AmazonEC2ReadOnlyAccess` policy to it. This policy allows read-only access to EC2 APIs.
1. When you launch the NGINX Plus instance, add this IAM role to the instance.

Alternatively, you can use the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environmental variables to provide credentials to nginx-asg-sync.

## Installation

To install nginx-asg-sync:

1. Download a software package for your OS with the latest version of nginx-asg-sync from the [Releases page](https://github.com/nginxinc/nginx-asg-sync/releases).
1. Install the package:

  For Amazon Linux or CentOS/RHEL, run: `$ sudo rpm -i <package-name>.rpm`

  For Ubuntu, run: `$ sudo dpkg -i <package-name>.deb`

## Configuration

As an example, we configure NGINX Plus to load balance two AWS Auto Scaling groups -- backend-group-one and backend-group-two. NGINX Plus routes requests to the appropriate Auto Scaling group based on the request URI:

* Requests for /backend-one go to Backend One group.
* Requests for /backend-two go to Backend Two group.

### NGINX Plus Configuration

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

   location /backend-two {
       proxy_set_header Host $host;
       proxy_pass http://backend-two;
   }
}

server {
    listen 8080;

    root /usr/share/nginx/html;

    location = / {
        return 302 /status.html;
    }

    location = /status.html {
    }

    location /status {
        access_log off;
        status;
    }

    location /upstream_conf {
        upstream_conf;
    }
}
```

* We declare two upstream groups – **backend-one** and **backend-two**, which correspond to our Auto Scaling groups. However, we do not add any servers to the upstream groups, because the servers will be added by nginx-aws-sync. The `state` directive names the file where the dynamically configurable list of servers is stored, enabling it to persist across restarts of NGINX Plus.
* We define a virtual server that listens on port 80. NGINX Plus passes requests for **/backend-one** to the instances of the Backend One group, and requests for **/backend-two** to the instances of the Backend Two group.
* We define a second virtual server listening on port 8080 and configure the NGINX Plus APIs on it, which are required by nginx-asg-sync:
  * The on-the-fly API is available at **127.0.0.1:8080/upstream_conf**
  * The status API is available at **127.0.0.1:8080/status**

### nginx-asg-sync Configuration

nginx-asg-sync is configured in the file **aws.yaml** in the **/etc/nginx** folder. For our example, we define the following configuration:

```yaml
region: us-west-2
upstream_conf_endpoint: http://127.0.0.1:8080/upstream_conf
status_endpoint: http://127.0.0.1:8080/status
sync_interval_in_seconds: 5
upstreams:
 - name: backend-one
   autoscaling_group: backend-one-group
   port: 80
   kind: http
 - name: backend-two
   autoscaling_group: backend-two-group
   port: 80
   kind: http
```

* The `region` key defines the AWS region where we deploy NGINX Plus and the Auto Scaling groups.
* The `upstream_conf` and the `status_endpoint` keys define the NGINX Plus API endpoints.
* The `sync_interval_in_seconds` key defines the synchronization interval: nginx-asg-sync checks for scaling updates every 5 seconds.
* The `upstreams` key defines the list of upstream groups. For each upstream group we specify:
  * `name` – The name we specified for the upstream block in the NGINX Plus configuration.
  * `autoscaling_group` – The name of the corresponding Auto Scaling group.
  * `port` – The port on which our backend applications are exposed.
  * `protocol` – The protocol of the traffic NGINX Plus load balances to the backend application, here `http`. If the application uses TCP/UDP, specify `stream` instead.

## Usage

nginx-asg-sync runs as a system service and supports the start/stop/restart commands.

For Ubuntu 14.04 and Amazon Linux, run: `$ sudo start|stop|restart nginx-asg-sync`

For Ubuntu 16.04 and CentOS7/RHEL7, run: `$ sudo service nginx-asg-sync start|stop|restart`

## Troubleshooting

If nginx-asg-sync doesn’t work as expected, check its log file available at **/var/log/nginx-aws-sync/nginx-aws-sync.log**.

## Building a Software Package

You can compile nginx-asg-sync and build a software package using the provided Makefile. Before you start building a package, make sure that the following software is installed on your system:
* make
* Docker

To build a software package, run: `$ make <os>`
where `<os>` is the target OS. The following values are allowed:
* `amazon` for Amazon Linux
* `centos7` for CentOS7/RHEL7
* `ubuntu-trusty` for Ubuntu 14.04
* `ubuntu-xenial` for Ubuntu 16.04

If you run make without any arguments, it will build software packages for all supported OSes.

## Support

Support from the [NGINX Professional Services Team](https://www.nginx.com/services/) is available when using nginx-asg-sync.
