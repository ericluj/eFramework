FROM golang:1.18-alpine as builder
WORKDIR /root
COPY ./  ./
RUN export GO111MODULE=on && CGO_ENABLED=0 GOOS=linux go build -o build/front /root/cmd/front/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/build/front ./
ENTRYPOINT ["/root/front"]