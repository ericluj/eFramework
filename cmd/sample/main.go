package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/ericluj/eFramework/common"
	"github.com/ericluj/eFramework/consul" // grpc使用consul做服务发现init
	"github.com/ericluj/eFramework/jaeger"
	"github.com/ericluj/eFramework/rpc/sample"
	log "github.com/ericluj/elog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serviceName = "sample"
	port        = 8000
)

func main() {
	log.Infof("start framework")

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
	// jaeger
	tracer, closer, err := jaeger.NewJaegerTracer(serviceName)
	defer closer.Close()
	if err != nil {
		log.Infof("NewJaegerTracer err: %v", err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(jaeger.ServerInterceptor(tracer)))
	sample.RegisterSampleServiceServer(server, &SampleService{})

	// 注册到consul
	consul.RegitserService(&consul.ConsulService{
		Name: serviceName,
		Tag:  []string{serviceName},
		IP:   common.GetLocalIP(),
		Port: port,
	})

	_ = server.Serve(ln)
}

func initHttpServer() {
	mux := runtime.NewServeMux()
	grpcServerEndpoint := fmt.Sprintf("localhost:%d", port)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := sample.RegisterSampleServiceHandlerFromEndpoint(context.Background(), mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Infof("http server error: %v", err)
		return
	}

	_ = http.ListenAndServe(":8001", mux)
}

type SampleService struct {
	sample.UnimplementedSampleServiceServer
}

func (s *SampleService) Health(ctx context.Context, in *sample.HealthRequest) (*sample.HealthResponse, error) {
	log.Infof("grpc health")
	return &sample.HealthResponse{Status: "ok"}, nil
}

func (s *SampleService) Search(ctx context.Context, in *sample.SearchRequest) (*sample.SearchResponse, error) {
	log.Infof("grpc search")
	return &sample.SearchResponse{Response: "haha"}, nil
}
