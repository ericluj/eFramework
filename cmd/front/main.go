package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"eFramework/common"
	"eFramework/consul" // grpc使用consul做服务发现init
	"eFramework/jaeger"
	"eFramework/rpc/front"
	"eFramework/rpc/sample"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serviceName = "front"
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
	// jaeger
	tracer, closer, err := jaeger.NewJaegerTracer(serviceName)
	defer closer.Close()
	if err != nil {
		fmt.Printf("NewJaegerTracer err: %v", err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(jaeger.ServerInterceptor(tracer)))
	front.RegisterFrontServiceServer(server, &FrontService{})

	// 注册到consul
	consul.RegitserService(&consul.ConsulService{
		Name: serviceName,
		Tag:  []string{serviceName},
		IP:   common.GetLocalIP(),
		Port: port,
	})

	server.Serve(ln)
}

func initHttpServer() {
	mux := runtime.NewServeMux()
	grpcServerEndpoint := fmt.Sprintf("localhost:%d", port)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := front.RegisterFrontServiceHandlerFromEndpoint(context.Background(), mux, grpcServerEndpoint, opts)
	if err != nil {
		fmt.Printf("http server error: %v", err)
		return
	}

	http.ListenAndServe(":8001", mux)
}

type FrontService struct {
	front.UnimplementedFrontServiceServer
}

func (s *FrontService) Health(ctx context.Context, in *front.HealthRequest) (*front.HealthResponse, error) {
	fmt.Println("grpc health")
	return &front.HealthResponse{Status: "ok"}, nil
}

func (s *FrontService) Sample(ctx context.Context, in *front.SampleRequest) (*front.SampleResponse, error) {
	fmt.Println("grpc sample")
	res := &front.SampleResponse{}
	data, err := GetSampleClient().Search(ctx, &sample.SearchRequest{Request: "front发出"})
	if err != nil {
		return res, err
	}
	res.Response = fmt.Sprintf("sample收到:%s", data.Response)
	return res, nil
}

var sampleOnce sync.Once
var sampleClient sample.SampleServiceClient

func GetSampleClient() sample.SampleServiceClient {
	sampleOnce.Do(func() {
		conn, err := common.NewClientConn("sample")
		if err != nil {
			panic(err)
		}
		sampleClient = sample.NewSampleServiceClient(conn)
	})
	return sampleClient
}
