## 0.2-1 (July 27, 2018)

IMPROVEMENTS:

* Add supporting documentation for the project https://github.com/nginxinc/nginx-asg-sync/pull/10
* Update package layout https://github.com/nginxinc/nginx-asg-sync/pull/9
* Use new NGINX Plus API https://github.com/nginxinc/nginx-asg-sync/pull/7

UPGRADE:

The upgrade process requires changing both NGINX Plus configuration and nginx-asg-sync configuration. Below are the recommended steps to follow:

1. Upgrade NGINX Plus to R14 or R15
2. Enable the new API in the NGINX Plus configuration while keeping the upstream_conf and the status API enabled. See an example of configuring the new API in the configuration section, but make sure to keep the upstream_conf and the status API.
3. Reload NGINX Plus to apply the updated configuration
4. Modify the /etc/nginx/aws.yaml file:
    * Remove the `upstream_conf_endpoint` and `status_endpoint` fields.
    * Add the `api_endpoint` field. See an example in the configuration section of the README.md
5. Download the Release 0.2 nginx-asg-sync package for your OS and upgrade the package using the OS tools (dpkg or rpm).
6. Check the logs of nginx-asg-sync to make sure that it is working properly after the upgrade.
7. Finally remove the upstream_conf and the status API from NGINX Plus configuration.
8. Reload NGINX Plus to apply the updated configuration

Note: the supported versions of NGINX Plus are R14 and higher.

## 0.1-2 (August 30, 2017)

IMPROVEMENTS:

* Make sure nginx-asg-sync works with NGINX Plus R13


## 0.1-1 (March 6, 2017)

Initial release
