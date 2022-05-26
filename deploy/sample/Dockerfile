FROM golang:1.18-alpine as builder
WORKDIR /root
COPY ./  ./
RUN export GO111MODULE=on && CGO_ENABLED=0 GOOS=linux go build -o build/sample /root/cmd/sample/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/build/sample ./
ENTRYPOINT ["/root/sample"]