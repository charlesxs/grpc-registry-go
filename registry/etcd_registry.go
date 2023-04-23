package registry

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
	"time"
)

type etcdRegistry struct {
	options *EtcdRegistryOptions

	c  *clientv3.Client
	em endpoints.Manager

	logger *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func newEtcdRegistry(client *clientv3.Client, options *EtcdRegistryOptions) (IRegistry, error) {
	em, err := endpoints.NewManager(client, options.AppName)
	if err != nil {
		return nil, err
	}

	r := &etcdRegistry{
		options: options,
		c:       client,
		em:      em,
		logger:  options.Logger,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	return r, nil
}

func (r *etcdRegistry) Register(addr string, port int, metadata interface{}) error {
	fqdn, endpointPath := r.makeRegisterPath(addr, port)
	ep := endpoints.Endpoint{Addr: fqdn, Metadata: metadata}

	// 注册自己
	err := r.register(endpointPath, ep)
	if err != nil {
		return err
	}

	// 后台goroutine 检查并续约
	go func() {
		ticker := time.NewTicker(r.options.Interval)
		for {
			select {
			case <-r.ctx.Done():
				r.logger.Sugar().Infof("registry stop register, endpoint=%s", fqdn)
				return
			case <-ticker.C:
			}

			err = r.register(endpointPath, ep)
			if err != nil {
				r.logger.Sugar().Error("registry register failed", zap.Error(err))
			}
		}
	}()

	return nil
}

func (r *etcdRegistry) Unregister(addr string, port int) error {
	_, endpointPath := r.makeRegisterPath(addr, port)

	// cancel
	r.cancel()

	// 重置cancel context
	r.ctx, r.cancel = context.WithCancel(context.Background())
	return r.em.DeleteEndpoint(r.ctx, endpointPath)
}

func (r *etcdRegistry) register(key string, ep endpoints.Endpoint) error {
	// minimum lease TTL is ttl-second
	resp, err := r.c.Grant(r.ctx, r.options.LeaseTTL)
	if err != nil {
		return fmt.Errorf("[%w] grant lease failed, error=%s", ErrRegistry, err.Error())
	}

	// 查看 key 是否存在, 存在则续租，不存在则创建并续租
	_, err = r.c.Get(r.ctx, key)
	if err == nil || err == rpctypes.ErrKeyNotFound {
		return r.em.AddEndpoint(r.ctx, key, ep, clientv3.WithLease(resp.ID))
	}

	return fmt.Errorf("[%w] register endpoint failed, error=%s, endpoint=%s", ErrRegistry, err.Error(), ep.Addr)
}

func (r *etcdRegistry) makeRegisterPath(addr string, port int) (string, string) {
	fqdn := fmt.Sprintf("%s:%d", addr, port)
	endpointPath := fmt.Sprintf("%s/%s", r.options.AppName, fqdn)
	return fqdn, endpointPath
}
