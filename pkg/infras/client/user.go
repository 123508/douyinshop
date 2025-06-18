package client

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/123508/douyinshop/kitex_gen/auth/authservice"
	"time"

	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/123508/douyinshop/kitex_gen/user/userservice"
	"github.com/123508/douyinshop/pkg/config"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var userClient userservice.Client

var authClientUser authservice.Client

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

	a, err := authservice.NewClient(
		config.Conf.UserConfig.ServiceName,                // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}

	authClientUser = a
}

func Register(ctx context.Context, req *user.RegisterReq) (bool, error) {
	_, err := userClient.Register(ctx, req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func Login(ctx context.Context, req *user.LoginReq) (string, error) {
	resp, err := userClient.Login(ctx, req)
	if err != nil {
		return "", err
	}

	tokenResp, err := userClient.DeliverTokenByRPC(ctx, &user.DeliverTokenReq{UserId: resp.UserId})

	if err != nil {
		return "", err
	}

	return tokenResp.Token, nil
}

func Delete(ctx context.Context, req *user.DeleteReq) (bool, error) {
	_, err := userClient.Delete(ctx, req)
	return err == nil, err
}

func GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	resp, err := userClient.GetUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Update(ctx context.Context, req *user.UpdateReq) (bool, error) {
	_, err := userClient.Update(ctx, req)
	return err == nil, err
}

func Logout(ctx context.Context, req *user.LogoutReq) (bool, error) {
	_, err := userClient.Logout(ctx, req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeliverToken(ctx context.Context, userId uint64) (string, error) {

	//用户角色身份获取
	roleResp, err := authClientUser.GetUserRoles(ctx, &auth.GetUserRolesReq{UserId: userId})

	if err != nil {
		return "", err
	}

	roleList := make([]string, 0, len(roleResp.Roles))

	for _, v := range roleResp.Roles {
		roleList = append(roleList, v.Code)
	}

	//用户权限身份获取
	permResp, err := authClientUser.GetUserPermissions(ctx, &auth.GetUserPermissionsReq{UserId: userId})

	if err != nil {
		return "", err
	}

	permList := make([]string, 0, len(permResp.Permissions))

	for _, v := range permResp.Permissions {
		permList = append(permList, v.Code)
	}

	resp, err := userClient.DeliverTokenByRPC(ctx, &user.DeliverTokenReq{
		UserId:    userId,
		RoleCodes: roleList,
		PermCodes: permList,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func VerifyToken(ctx context.Context, req *user.VerifyTokenReq) (uint64, string, error) {
	resp, err := userClient.VerifyTokenByRPC(ctx, req)
	if err != nil || resp == nil {
		return 0, "", err
	}
	return resp.UserId, req.Token, nil
}

func ListUsers(ctx context.Context, req *user.ListUsersReq) (*user.ListUsersResp, error) {
	resp, err := userClient.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ChangePassword(ctx context.Context, req *user.ChangePasswordReq) (bool, error) {
	_, err := userClient.ChangePassword(ctx, req)
	return err == nil, err
}

func ForgotPassword(ctx context.Context, req *user.ForgotPasswordReq) (bool, error) {
	_, err := userClient.ForgotPassword(ctx, req)
	return err == nil, err
}

func ResetPassword(ctx context.Context, req *user.ResetPasswordReq) (bool, error) {
	_, err := userClient.ResetPassword(ctx, req)
	return err == nil, err
}

func BindEmail(ctx context.Context, req *user.BindEmailReq) (bool, error) {
	_, err := userClient.BindEmail(ctx, req)
	return err == nil, err
}

func UnbindEmail(ctx context.Context, req *user.UnbindEmailReq) (bool, error) {
	_, err := userClient.UnbindEmail(ctx, req)
	return err == nil, err
}

func FreezeUser(ctx context.Context, req *user.FreezeUserReq) (bool, error) {
	_, err := userClient.FreezeUser(ctx, req)
	return err == nil, err
}

func UnfreezeUser(ctx context.Context, req *user.UnfreezeUserReq) (bool, error) {
	_, err := userClient.UnfreezeUser(ctx, req)
	return err == nil, err
}
