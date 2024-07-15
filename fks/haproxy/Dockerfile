FROM perconalab/percona-server-mysql-operator:main-haproxy
COPY --from=flyio/percona-server-mysql-operator:js-fks /opt/percona-server-mysql-operator /opt/percona

# FKS: Create this directory until we support emptyDir mounts in FKS
RUN mkdir -p /etc/haproxy/mysql
