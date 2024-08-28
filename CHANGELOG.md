# Changelog

Starting with version 1.0.0 an automatically generated list of changes can be found on the [GitHub Releases page](https://github.com/nginxinc/nginx-asg-sync/releases).

## 0.5.0 (February 24, 2021)

IMPROVEMENTS:

- Add InService option for AWS [#39](https://github.com/nginxinc/nginx-asg-sync/pull/39)
- Update log format [#42](https://github.com/nginxinc/nginx-asg-sync/pull/42)

## 0.4-1 (November 22, 2019)

IMPROVEMENTS:

- Add support to set upstream server parameters: `max_conns`, `max_fails`, `fail_timeout` and `slow_start` in the
  configuration file. <https://github.com/nginxinc/nginx-asg-sync/pull/33>
- Add support to use wildcards in the names of AWS Auto Scaling groups.
  <https://github.com/nginxinc/nginx-asg-sync/pull/29/>
- Allow nginx-asg-sync to detect the region where it is running (use `region: self` in the configuration file).
  <https://github.com/nginxinc/nginx-asg-sync/pull/27>

## 0.3-1 (September 4, 2019)

IMPROVEMENTS:

- Add support for Azure Virtual Machine Scale Sets <https://github.com/nginxinc/nginx-asg-sync/pull/24>
- Create separate documentation for the configuration for different cloud providers: [aws](examples/aws.md) and
  [azure](examples/azure.md).
- Ubuntu 14.04 (Trusty) is no longer supported.

UPGRADE:

The upgrade process requires changing the configuration file name. Below are the recommended steps to follow:

1. Change the name of the configuration file from `/etc/nginx/aws.yaml` to `/etc/nginx/config.yaml`.
2. Download the Release 0.3 nginx-asg-sync package for your OS and upgrade the package using the OS tools (dpkg or rpm).
3. Check the logs of nginx-asg-sync to make sure that it is working properly after the upgrade.

Note: the supported versions of NGINX Plus are R18 and higher.

## 0.2-1 (July 27, 2018)

IMPROVEMENTS:

- Add supporting documentation for the project <https://github.com/nginxinc/nginx-asg-sync/pull/10>
- Update package layout <https://github.com/nginxinc/nginx-asg-sync/pull/9>
- Use new NGINX Plus API <https://github.com/nginxinc/nginx-asg-sync/pull/7>

UPGRADE:

The upgrade process requires changing both NGINX Plus configuration and nginx-asg-sync configuration. Below are the
recommended steps to follow:

1. Upgrade NGINX Plus to R14 or R15
2. Enable the new API in the NGINX Plus configuration while keeping the upstream_conf and the status API enabled. See an
   example of configuring the new API in the configuration section, but make sure to keep the upstream_conf and the
   status API.
3. Reload NGINX Plus to apply the updated configuration
4. Modify the /etc/nginx/aws.yaml file:
   - Remove the `upstream_conf_endpoint` and `status_endpoint` fields.
   - Add the `api_endpoint` field. See an example in the configuration section of the README.md
5. Download the Release 0.2 nginx-asg-sync package for your OS and upgrade the package using the OS tools (dpkg or rpm).
6. Check the logs of nginx-asg-sync to make sure that it is working properly after the upgrade.
7. Finally remove the upstream_conf and the status API from NGINX Plus configuration.
8. Reload NGINX Plus to apply the updated configuration

Note: the supported versions of NGINX Plus are R14 and higher.

## 0.1-2 (August 30, 2017)

IMPROVEMENTS:

- Make sure nginx-asg-sync works with NGINX Plus R13

## 0.1-1 (March 6, 2017)

Initial release
