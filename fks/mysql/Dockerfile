FROM perconalab/percona-server-mysql-operator:main-psmysql
COPY --from=flyio/percona-server-mysql-operator:fks-22 /opt/percona-server-mysql-operator /opt/percona
ENV PT_VERSION 3.5.7-1.el9
USER root
RUN dnf update && dnf install microdnf && percona-release enable tools release && microdnf -y install percona-toolkit-${FULL_PERCONA_VERSION}
USER 2
