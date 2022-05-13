package test

import (
	"context"
	"testing"
	"time"

	"eFramework/consul"
	"eFramework/rpc/sample"

	"google.golang.org/grpc"
)

func TestConsulRegister(t *testing.T) {
	consul.RegitserService(&consul.ConsulService{
		Name: "testregister",
		Tag:  []string{"testregister"},
		IP:   "127.0.0.1",
		Port: 50051,
	})
}

func TestConsulFind(t *testing.T) {
	target := "consul://127.0.0.1:8500/sample"
	conn, err := grpc.DialContext(context.Background(), target, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		t.Errorf("DialContext error: %v", err)
	}
	defer conn.Close()
	c := sample.NewSampleServiceClient(conn)

	resp, err := c.Search(context.Background(), &sample.SearchRequest{Request: "TestConsulFind"})
	if err != nil {
		t.Errorf("Search error: %v", err)
	}
	t.Logf("resp: %s", resp.GetResponse())
	time.Sleep(time.Second * 2)
}
