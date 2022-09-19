package consul

import (
	"context"
	"fmt"
	"strings"
	"sync"

	log "github.com/ericluj/elog"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const scheme = "consul"

func init() {
	resolver.Register(&consulBuilder{})
}

type consulBuilder struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	ctx, cancel := context.WithCancel(context.Background())
	cr := &consulResolver{
		ctx:       ctx,
		cancel:    cancel,
		lastIndex: 0,
		address:   target.URL.Host,
		name:      strings.Trim(target.URL.Path, "/"),
		cc:        cc,
	}

	cr.wg.Add(1)
	go cr.watcher()

	return cr, nil
}

func (*consulBuilder) Scheme() string {
	return scheme
}

type consulResolver struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	address   string
	name      string
	lastIndex uint64
	cc        resolver.ClientConn
}

func (cr *consulResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (cr *consulResolver) Close() {}

func (cr *consulResolver) watcher() {
	config := api.DefaultConfig()
	config.Address = cr.address

	client, err := api.NewClient(config)
	if err != nil {
		log.Infof("consul client error: %v", err)
		return
	}

	for {
		services, metainfo, err := client.Health().Service(cr.name, cr.name, true, &api.QueryOptions{WaitIndex: cr.lastIndex})
		if err != nil {
			log.Infof("service from consul errror: %v", err)
		}
		cr.lastIndex = metainfo.LastIndex

		var addrs []resolver.Address
		for _, s := range services {
			addr := fmt.Sprintf("%v:%v", s.Service.Address, s.Service.Port)
			addrs = append(addrs, resolver.Address{Addr: addr})
		}
		cr.cc.NewAddress(addrs)
		cr.cc.NewServiceConfig(cr.name)
	}
}
