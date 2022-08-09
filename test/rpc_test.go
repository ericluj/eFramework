package test

import (
	"context"
	"testing"

	"github.com/ericluj/eFramework/rpc/sample"
	"google.golang.org/grpc"
)

func TestRpc(t *testing.T) {
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	if err != nil {
		t.Errorf("grpc.Dial err: %v", err)
		return
	}
	defer conn.Close()
	client := sample.NewSampleServiceClient(conn)
	resp, err := client.Search(context.Background(), &sample.SearchRequest{
		Request: "hello",
	})
	if err != nil {
		t.Errorf("client.Search err: %v", err)
		return
	}

	t.Logf("resp: %s", resp.GetResponse())
}
