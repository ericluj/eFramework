package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"eFramework/consul" // grpc使用consul做服务发现init
	"eFramework/rpc/sample"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serviceName = "sample"
	port        = 8000
)

func main() {
	fmt.Println("start framework")

	ln, err := listen(port)
	if err != nil {
		panic(err)
	}

	// http server
	go initHttpServer()

	// rpc server
	go initGrpcServer(ln)

	select {}
}

func listen(port int) (ln net.Listener, err error) {
	ln, err = net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}
	return ln, nil
}

func initGrpcServer(ln net.Listener) {
	server := grpc.NewServer()
	sample.RegisterSampleServiceServer(server, &SampleService{})

	// 注册到consul
	consul.RegitserService(&consul.ConsulService{
		Name: serviceName,
		Tag:  []string{serviceName},
		IP:   "127.0.0.1", // TODO: 需要改为获取本机ip
		Port: port,
	})

	server.Serve(ln)
}

func initHttpServer() {
	mux := runtime.NewServeMux()
	grpcServerEndpoint := fmt.Sprintf("localhost:%d", port)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := sample.RegisterSampleServiceHandlerFromEndpoint(context.Background(), mux, grpcServerEndpoint, opts)
	if err != nil {
		fmt.Printf("http server error: %v", err)
		return
	}

	http.ListenAndServe(":8001", mux)
}

type SampleService struct {
	sample.UnimplementedSampleServiceServer
}

func (s *SampleService) Health(ctx context.Context, in *sample.HealthRequest) (*sample.HealthResponse, error) {
	fmt.Println("grpc health")
	return &sample.HealthResponse{Status: "ok"}, nil
}

func (s *SampleService) Search(ctx context.Context, in *sample.SearchRequest) (*sample.SearchResponse, error) {
	fmt.Println("grpc search")
	return &sample.SearchResponse{Response: "haha"}, nil
}
