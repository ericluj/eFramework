package common

import (
	"context"
	"eFramework/consul"
	"fmt"
	"net"

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
	target := fmt.Sprintf("consul://%s/%s", consul.ConsulAddress, serviceName)
	conn, err := grpc.DialContext(context.Background(), target, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	return conn, err
}
