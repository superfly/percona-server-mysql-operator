FROM perconalab/percona-server-mysql-operator:main-psmysql

ENV PT_VERSION 3.5.7-1.el9
USER root
RUN dnf update && dnf install microdnf && percona-release enable tools release && microdnf -y install percona-toolkit-${PT_VERSION}
USER mysql

COPY --from=flyio/percona-server-mysql-operator:fks-23 /opt/percona-server-mysql-operator /opt/percona
