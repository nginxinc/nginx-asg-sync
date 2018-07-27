## 0.2-1 (July 27, 2018)

IMPROVEMENTS:

* Add supporting documentation for the project https://github.com/nginxinc/nginx-asg-sync/pull/10
* Update package layout https://github.com/nginxinc/nginx-asg-sync/pull/9
* Use new NGINX Plus API https://github.com/nginxinc/nginx-asg-sync/pull/7

UPGRADE:

* Remove the previous version of nginx-asg-sync e.g. `dpkg --remove nginx-asg-sync`
* Deploy the new version to your NGINX Plus instance in AWS and install it `dpkg -i nginx-asg-sync_0.2-1-xenial_amd64.deb`
* Update the `/etc/nginx/aws.yaml` to the new format (example in the [configuration section](https://github.com/nginxinc/nginx-asg-sync#nginx-asg-sync-configuration) of the README.md))
* Reload NGINX Plus

Note: the supported versions of NGINX Plus are R14 and higher.

## 0.1-2 (August 30, 2017)

IMPROVEMENTS:

* Make sure nginx-asg-sync works with NGINX Plus R13


## 0.1-1 (March 6, 2017)

Initial release
