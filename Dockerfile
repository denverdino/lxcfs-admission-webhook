FROM alpine:latest

ADD lxcfs-admission-webhook /lxcfs-admission-webhook
ENTRYPOINT ["./lxcfs-admission-webhook"]