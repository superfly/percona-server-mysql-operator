FROM redhat/ubi9-minimal as users

RUN microdnf install shadow-utils jq tar -y

RUN groupadd -g 1002 pmm-agent && \
    useradd -u 1002 -r -g pmm-agent -s /sbin/nologin \
    -d /usr/local/percona/pmm2 \
    -c "PMM 2.X Client User" pmm-agent


FROM flyio/ps8

# Merge XtraBackub, PMM agent and Percona Toolkit installation into the main image for sidecar support
# https://github.com/percona/percona-docker/blob/main/percona-xtrabackup-8.0/Dockerfile#L58-L60
# https://github.com/percona/percona-docker/blob/main/percona-toolkit/Dockerfile#L14-L15

ENV XTRABACKUP_VERSION 8.0.35-31.1
ENV PS_VERSION 8.0.36-28.1
ENV OS_VER el9
ENV FULL_PERCONA_VERSION "$PS_VERSION.$OS_VER"
ENV FULL_PERCONA_XTRABACKUP_VERSION "$XTRABACKUP_VERSION.$OS_VER"

ENV PATH /usr/local/percona/pmm2/bin/:$PATH

ENV PT_VERSION 3.5.7-1.el9

USER root
RUN dnf update && dnf install microdnf && percona-release enable tools release && \
    curl -Lf -o /tmp/libev.rpm https://downloads.percona.com/downloads/packaging/libev-4.33-5.el9.x86_64.rpm && \
    rpm -i /tmp/libev.rpm && \
    rm -rf /tmp/libev.rpm && \
    microdnf -y install \
    percona-toolkit-${PT_VERSION} \
    percona-xtrabackup-80-${FULL_PERCONA_XTRABACKUP_VERSION} \
    percona-server-shared-${FULL_PERCONA_VERSION} \
    percona-server-client-${FULL_PERCONA_VERSION} \
    pmm2-client \
    socat procps-ng qpress lz4 shadow-utils jq tar zstd && \
    microdnf clean all && \
    rm -rf /var/cache/dnf /var/cache/yum

# Create directories created as emptyDir volumes in the operator
RUN install -d -o 1001 -g 0 -m 0775 /backup
RUN install -d -o 1001 -g 0 -m 0775 /var/log/xtrabackup

COPY --from=flyio/percona-server-mysql-operator:js-fks /opt/percona-server-mysql-operator /opt/percona

# Give access to pmm agent files
RUN chown mysql:mysql -R /usr/local/percona/pmm2 /etc/my.cnf.d  && chmod 755 /opt/percona/pmm-prerun.sh
USER mysql
