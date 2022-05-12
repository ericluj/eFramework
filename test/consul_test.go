package test

import (
	"testing"

	"eFramework/consul"
)

func TestConsulRegister(t *testing.T) {
	consul.RegitserService(&consul.ConsulService{
		Name: "testregister",
		Tag:  []string{"testregister"},
		IP:   "127.0.0.1",
		Port: 50051,
	})
}
