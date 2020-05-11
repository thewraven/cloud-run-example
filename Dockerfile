FROM golang:1.14-alpine as builder
WORKDIR /opt
COPY . /opt
RUN go build .

FROM alpine:3.11
RUN addgroup -S webgroup && adduser -S appgroup -G webgroup
USER appgroup
WORKDIR /opt
COPY --from=builder /opt/cloud-run-example .
ENTRYPOINT [ "/opt/cloud-run-example" ]