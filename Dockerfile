FROM golang:1.13

ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN mkdir /frp-manager
ADD . /frp-manager
WORKDIR /frp-manager
RUN go build -o /bin/manager ./cmd

FROM alpine:3.10

COPY --from=0 /frp-manager/bin/manager .

ENTRYPOINT ["/manager"]