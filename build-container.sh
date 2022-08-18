#!/bin/bash
docker build -t registry.cn-hangzhou.aliyuncs.com/denverdino/lxcfs:5.0.2 lxcfs-image
docker push registry.cn-hangzhou.aliyuncs.com/denverdino/lxcfs:5.0.2

./build
