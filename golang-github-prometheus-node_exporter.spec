
%global debug_package   %{nil}

%global provider        github
%global provider_tld    com
%global project         prometheus
%global repo            node_exporter
# https://github.com/prometheus/node_exporter
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path     %{provider_prefix}
%global commit          052cee4ae2677a7cb42542615f6c858e35d6e9e3
%global shortcommit     %(c=%{commit}; echo ${c:0:7})
%global gopathdir       %{_sourcedir}/go
%global upstream_ver    0.15.2
%global rpm_ver         %(v=%{upstream_ver}; echo ${v//-/_})
%global download_prefix %{provider}.%{provider_tld}/openshift/%{repo}

Name:		golang-%{provider}-%{project}-%{repo}
Version:	%{rpm_ver}
Release:	1.git%{shortcommit}%{?dist}
Summary:	Prometheus exporter for hardware and OS metrics exposed by *NIX kernels
License:	ASL 2.0
URL:		https://prometheus.io/
Source0:	https://%{download_prefix}/archive/%{commit}/%{repo}-%{shortcommit}.tar.gz

# e.g. el6 has ppc64 arch without gcc-go, so EA tag is required
# promu-based packages FTBFS on aarch64 (#1487462)
ExclusiveArch:  %{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 %{arm} ppc64le s390x}
# If go_compiler is not set to 1, there is no virtual provide. Use golang instead.
BuildRequires: %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
BuildRequires: glibc-static
BuildRequires: prometheus-promu

%description
%{summary}

%package -n %{project}-%{repo}
Summary:        %{summary}
Provides:       prometheus-node_exporter = %{version}-%{release}

%description -n %{project}-%{repo}
%{summary}

%prep
%setup -q -n %{repo}-%{commit}

%build
# Go expects a full path to the sources which is not included in the source
# tarball so create a link with the expected path
mkdir -p %{gopathdir}/src/%{provider}.%{provider_tld}/%{project}
ln -s `pwd` %{gopathdir}/src/%{import_path}
export GOPATH=%{gopathdir}

make build BUILD_PROMU=false

%install
install -d %{buildroot}%{_bindir}
install -D -p -m 0755 node_exporter %{buildroot}%{_bindir}/node_exporter

%files -n %{project}-%{repo}
%license LICENSE NOTICE
%doc CHANGELOG.md CONTRIBUTING.md MAINTAINERS.md README.md
%{_bindir}/node_exporter

%changelog
* Wed Jan 17 2018 Paul Gier <pgier@redhat.com> - 0.15.2-1
- upgrade to 0.15.2

* Wed Jan 10 2018 Yaakov Selkowitz <yselkowi@redhat.com> - 0.15.1-2
- Rebuilt for ppc64le, s390x enablement

* Thu Nov 09 2017 Paul Gier <pgier@redhat.com> - 0.15.1-1
- upgrade to 0.15.1
- added build requires on glibc-static

* Wed Aug 30 2017 Paul Gier <pgier@redhat.com> - 0.14.0-1
- Initial package creation


