package client

import (
	"context"

	"github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/123508/douyinshop/kitex_gen/auth/authservice"
	"github.com/123508/douyinshop/pkg/config"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var authClient authservice.Client

func initAuthRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := authservice.NewClient(
		config.Conf.AuthConfig.ServiceName,                // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	authClient = c
}

func DeliverToken(ctx context.Context, req *auth.DeliverTokenReq) (string, error) {
	resp, err := authClient.DeliverTokenByRPC(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func VerifyToken(ctx context.Context, req *auth.VerifyTokenReq) (uint32, string, error) {
	resp, err := authClient.VerifyTokenByRPC(ctx, req)
	if err != nil {
		return 0, "", err
	}
	if resp.Res == false {
		return 0, "", nil
	}
	return resp.UserId, req.Token, nil
}
