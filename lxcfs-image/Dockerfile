FROM centos:7 as yesbuilder
RUN yum -y update
RUN yum -y install fuse-devel pam-devel wget install gcc automake autoconf libtool make cmake help2man && \
    yum install --assumeyes python3-pip && \
    pip3 install scikit-build && \
    pip3 install jinja2 meson Ninja 
ENV LXCFS_VERSION 5.0.2
RUN wget https://linuxcontainers.org/downloads/lxcfs/lxcfs-$LXCFS_VERSION.tar.gz && \
    mkdir /lxcfs && \
    tar xzf lxcfs-$LXCFS_VERSION.tar.gz -C / && \
    cd /lxcfs-5.0.2/ && \
    meson setup -Dinit-script=systemd --prefix=/usr /lxcfs && \
    meson compile -C /lxcfs

FROM centos:7
STOPSIGNAL SIGINT
COPY --from=yesbuilder /lxcfs/lxcfs                    /lxcfs/lxcfs
COPY --from=yesbuilder /lxcfs/liblxcfs.so              /lxcfs/liblxcfs.so
COPY --from=yesbuilder /usr/lib64/libfuse.so.2.9.2     /lxcfs/libfuse.so.2.9.2
COPY --from=yesbuilder /usr/lib64/libulockmgr.so.1.0.1 /lxcfs/libulockmgr.so.1.0.1

COPY start.sh /
CMD ["/start.sh"]

