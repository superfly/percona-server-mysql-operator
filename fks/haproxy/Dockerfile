FROM perconalab/percona-server-mysql-operator:main-haproxy
COPY --from=flyio/percona-server-mysql-operator:js-fks /opt/percona-server-mysql-operator /opt/percona

# FKS: Create this directory until we support emptyDir mounts in FKS
USER root
RUN install -d -o 1001 -g 0 -m 0775 /etc/haproxy/mysql
USER mysql
