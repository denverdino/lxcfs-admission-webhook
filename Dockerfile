FROM golang:1.18.4-alpine as yesbuilder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR $GOPATH/src/github.com/lxcfs-admission-webhook
ADD . $GOPATH/src/github.com/lxcfs-admission-webhook
RUN go build . && \
    mv ./lxcfs-admission-webhook /usr/local/bin/

FROM alpine:latest
COPY --from=yesbuilder /usr/local/bin/lxcfs-admission-webhook /usr/local/bin/lxcfs-admission-webhook
ENTRYPOINT [ "/usr/local/bin/lxcfs-admission-webhook" ]
