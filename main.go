package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"eFramework/rpc/sample"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("start framework")

	port := 8000

	ln, err := listen(port)
	if err != nil {
		panic(err)
	}

	// rpc server
	go initGrpcServer(ln)

	// http server
	go initHttpServer(ln)

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

	server.Serve(ln)
}

func initHttpServer(ln net.Listener) {
	mux := http.NewServeMux()
	mux.Handle("/search", http.HandlerFunc(searchHandler))
	http.Serve(ln, mux)
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
