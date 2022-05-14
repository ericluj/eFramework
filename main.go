package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"eFramework/consul" // grpc使用consul做服务发现init
	"eFramework/rpc/sample"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
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

	// TODO: 后面cmux可以自己封装学习一下
	m := cmux.New(ln)
	grpcLn := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpLn := m.Match(cmux.Any()) // 除了grpc的剩下都当作http使用

	// rpc server
	go initGrpcServer(grpcLn)

	// http server
	go initHttpServer(httpLn)

	if err := m.Serve(); err != nil {
		panic(err)
	}

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

func initHttpServer(ln net.Listener) {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(healthHandler))
	http.Serve(ln, mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

type SampleService struct {
	sample.UnimplementedSampleServiceServer
}

func (s *SampleService) Search(ctx context.Context, in *sample.SearchRequest) (*sample.SearchResponse, error) {
	fmt.Println("grpc search")
	return &sample.SearchResponse{Response: "haha"}, nil
}
