package common

import (
	"context"
	"eFramework/consul"
	"eFramework/jaeger"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func NewClientConn(serviceName string) (grpc.ClientConnInterface, error) {
	// jaeger
	tracer, closer, err := jaeger.NewJaegerTracer(serviceName)
	defer closer.Close()
	if err != nil {
		log.Infof("NewJaegerTracer err: %v", err)
	}

	target := fmt.Sprintf("consul://%s/%s", consul.ConsulAddress, serviceName)
	conn, err := grpc.DialContext(
		context.Background(),
		target,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(jaeger.ClientInterceptor(tracer)),
	)
	return conn, err
}
