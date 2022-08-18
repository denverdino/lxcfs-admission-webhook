#!/bin/bash

# Cleanup
nsenter -m/proc/1/ns/mnt fusermount -u /var/lib/lxcfs 2> /dev/null || true
nsenter -m/proc/1/ns/mnt [ -L /etc/mtab ] || \
        sed -i "/^lxcfs \/var\/lib\/lxcfs fuse.lxcfs/d" /etc/mtab

# remove /var/lib/lxcfs
rm -rf /var/lib/lxcfs/*

# Prepare
mkdir -p /usr/local/lib/lxcfs /var/lib/lxcfs /usr/lib64/lxcfs/

# Update lxcfs
cp -f /lxcfs/lxcfs /usr/local/bin/lxcfs
cp -f /lxcfs/liblxcfs.so /usr/local/lib/lxcfs/liblxcfs.so
cp -f /lxcfs/liblxcfs.so /usr/lib64/lxcfs/liblxcfs.so

cp -f /lxcfs/libfuse.so.2.9.2     /usr/lib64/libfuse.so.2.9.2
cp -f /lxcfs/libulockmgr.so.1.0.1 /usr/lib64/libulockmgr.so.1.0.1

rm -f /usr/lib64/libfuse.so.2 /usr/lib64/libulockmgr.so.1

ln -s /usr/lib64/libfuse.so.2.9.2     /usr/lib64/libfuse.so.2
ln -s /usr/lib64/libulockmgr.so.1.0.1 /usr/lib64/libulockmgr.so.1

# Mount
exec nsenter -m/proc/1/ns/mnt /usr/local/bin/lxcfs /var/lib/lxcfs/

