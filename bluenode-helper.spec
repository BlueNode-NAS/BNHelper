Name:           bluenode-helper
Version:        1.0.0
Release:        1%{?dist}
Summary:        Backend for BlueNode Server OS

License:        Proprietary
URL:            https://github.com/BlueNode-NAS/BNHelper
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang
Requires:       systemd

%description
BlueNode Helper is a backend service for BlueNode NAS OS that provides
HTTP API functionality over Unix domain sockets.

%prep
%setup -q

%build
go build -v -ldflags "-X main.Version=%{version} -X main.BuildDate=$(date -u +%%Y-%%m-%%dT%%H:%%M:%%SZ) -X main.GitCommit=%{?git_commit}" -o %{name}

%install
rm -rf $RPM_BUILD_ROOT
install -d $RPM_BUILD_ROOT%{_bindir}
install -m 0755 %{name} $RPM_BUILD_ROOT%{_bindir}/%{name}

install -d $RPM_BUILD_ROOT%{_unitdir}
cat > $RPM_BUILD_ROOT%{_unitdir}/%{name}.service <<EOF
[Unit]
Description=BlueNode Helper Service
After=network.target

[Service]
Type=simple
ExecStart=%{_bindir}/%{name}
Restart=on-failure
RestartSec=5s
User=root
Group=root

[Install]
WantedBy=multi-user.target
EOF

%files
%{_bindir}/%{name}
%{_unitdir}/%{name}.service

%changelog
* Thu Jan 01 2026 BlueNode-NAS Team - 1.0.0-1
- Initial RPM package release
