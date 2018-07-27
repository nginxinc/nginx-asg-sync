# distribution specific definitions
%define use_systemd (0%{?rhel} && 0%{?rhel} >= 7)

Summary: NGINX Plus integration with AWS Auto Scaling groups
Name: nginx-asg-sync
Version: 0.2
Release: 2%{?dist}
Vendor: Nginx Software, Inc.
URL: https://github.com/nginxinc/nginx-asg-sync
Packager: Nginx Software, Inc. <https://www.nginx.com>

Source0: aws.yaml.example
Source1: COPYRIGHT
Source2: nginx-asg-sync.conf
Source3: nginx-asg-sync.logrotate
Source4: nginx-asg-sync.service

License: 2-clause BSD-like license
Group: System Environment/Daemons

%if %{use_systemd}
Requires: systemd
%else
Requires: upstart >= 0.6.5
%endif

%description
This package contains software that integrates NGINX Plus
with AWS Auto Scaling groups

%prep

%build

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{_sbindir}/

install -m 755 -p /build_output/nginx-asg-sync %{buildroot}%{_sbindir}/

%if %{use_systemd}
mkdir -p %{buildroot}%{_unitdir}
install -m 644 %{SOURCE4}  %{buildroot}%{_unitdir}/nginx-asg-sync.service
%else
mkdir -p %{buildroot}%{_sysconfdir}/init
install -m 644 %{SOURCE2} %{buildroot}%{_sysconfdir}/init/nginx-asg-sync.conf
%endif

mkdir -p %{buildroot}%{_sysconfdir}/logrotate.d
install -m 644 -p %{SOURCE3} %{buildroot}%{_sysconfdir}/logrotate.d/nginx-asg-sync

mkdir -p %{buildroot}%{_datadir}/doc/nginx-asg-sync
install -m 644 %{SOURCE1} -p %{buildroot}%{_datadir}/doc/nginx-asg-sync/

mkdir -p %{buildroot}%{_sysconfdir}/nginx
install -m 644 %{SOURCE0} -p %{buildroot}%{_sysconfdir}/nginx/

%{__mkdir} -p %{buildroot}%{_localstatedir}/log/nginx-asg-sync

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root)

%{_sbindir}/nginx-asg-sync

%if %{use_systemd}
%{_unitdir}/nginx-asg-sync.service
%else
%{_sysconfdir}/init/nginx-asg-sync.conf
%endif

%dir %{_sysconfdir}/nginx
%config(noreplace) %{_sysconfdir}/nginx/aws.yaml.example
%config(noreplace) %{_sysconfdir}/logrotate.d/nginx-asg-sync

%dir %{_datadir}/doc/nginx-asg-sync
%{_datadir}/doc/nginx-asg-sync/*

%attr(0755,nginx,nginx) %dir %{_localstatedir}/log/nginx-asg-sync

%post
%if %{use_systemd}
if [ $1 -eq 1 ]; then
	/usr/bin/systemctl preset nginx-asg-sync.service >/dev/null 2>&1 ||:
fi
%endif

%preun
if [ $1 -eq 0 ]; then
%if %use_systemd
	/usr/bin/systemctl --no-reload disable nginx-asg-sync.service >/dev/null 2>&1 ||:
	/usr/bin/systemctl stop nginx-asg-sync.service >/dev/null 2>&1 ||:
%else
	/sbin/stop nginx-asg-sync >/dev/null 2>&1 ||:
%endif
fi

%postun
%if %use_systemd
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:
%endif
if [ $1 -ge 1 ]; then
%if %use_systemd
	/usr/bin/systemctl restart nginx-asg-sync.service >/dev/null 2>&1 ||:
%else
	/sbin/restart nginx-asg-sync >/dev/null 2>&1 ||:
%endif
fi

%changelog
* Fri Jul 27 2018 Peter Kelly <peter.kelly@nginx.com>
- 0.2-1
- Add supporting guides for contributing and changelog
- Update package layout
- Use new NGINX Plus API

* Wed Aug 30 2017 Michael Pleshakov <michael@nginx.com>
- 0.1-2
- Make sure nginx-asg-sync works with NGINX Plus R13

* Fri Mar 03 2017 Michael Pleshakov <michael@nginx.com>
- 0.1-1
- First release
