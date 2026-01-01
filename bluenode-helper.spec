Name:           bluenode-helper
Version:        1.0.0
Release:        1%{?dist}
Summary:        Backend for BlueNode Server OS

License:        Proprietary
URL:            https://github.com/BlueNode-NAS/BNHelper
Source0:        %{name}
Source1:        %{name}.service

Requires:       systemd

%description
BlueNode Helper is a backend service for BlueNode NAS OS that provides
HTTP API functionality over Unix domain sockets.

%prep
# No prep needed - using pre-built binary

%build
# No build needed - using pre-built binary

%install
rm -rf $RPM_BUILD_ROOT
install -d $RPM_BUILD_ROOT%{_bindir}
install -m 0755 %{SOURCE0} $RPM_BUILD_ROOT%{_bindir}/%{name}

install -d $RPM_BUILD_ROOT%{_unitdir}
install -m 0644 %{SOURCE1} $RPM_BUILD_ROOT%{_unitdir}/%{name}.service

%files
%{_bindir}/%{name}
%{_unitdir}/%{name}.service

%changelog
* Thu Jan 01 2026 BlueNode-NAS Team - 1.0.0-1
- Initial RPM package release
