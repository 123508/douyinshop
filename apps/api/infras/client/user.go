package client

import (
	"context"
	"time"

	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/123508/douyinshop/kitex_gen/user/userservice"
	"github.com/123508/douyinshop/pkg/config"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var userClient userservice.Client

func initUserRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := userservice.NewClient(
		config.Conf.UserConfig.ServiceName,                // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func Register(ctx context.Context, req *user.RegisterReq) (uint32, error) {
	resp, err := userClient.Register(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.UserId, nil
}

func Login(ctx context.Context, req *user.LoginReq) (uint32, error) {
	resp, err := userClient.Login(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.UserId, nil
}
