package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"eFramework/consul"
	_ "eFramework/consul" // grpc使用consul做服务发现init
	"eFramework/rpc/sample"

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

	// http server
	go initHttpServer(ln)

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

func initHttpServer(ln net.Listener) {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(healthHandler))
	mux.Handle("/search", http.HandlerFunc(searchHandler))
	http.Serve(ln, mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("http search")
	fmt.Fprintf(w, "http search")
}

type SampleService struct {
	sample.UnimplementedSampleServiceServer
}

func (s *SampleService) Search(ctx context.Context, in *sample.SearchRequest) (*sample.SearchResponse, error) {
	fmt.Println("grpc search")
	return &sample.SearchResponse{Response: "haha"}, nil
}
