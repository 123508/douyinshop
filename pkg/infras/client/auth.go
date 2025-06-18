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

// AssignRoleToUser 给用户分配角色
func AssignRoleToUser(ctx context.Context, userId, roleId uint64) (bool, error) {

	req := &auth.AssignRoleToUserReq{
		UserId: userId,
		RoleId: roleId,
	}
	_, err := authClient.AssignRoleToUser(ctx, req)
	return err == nil, err
}

// RemoveRoleFromUser 移除用户的角色
func RemoveRoleFromUser(ctx context.Context, userId, roleId uint64) (bool, error) {
	req := &auth.RemoveRoleFromUserReq{
		UserId: userId,
		RoleId: roleId,
	}
	_, err := authClient.RemoveRoleFromUser(ctx, req)
	return err == nil, err
}

// GetUserRoles 查询用户所有角色
func GetUserRoles(ctx context.Context, userId uint64) (*auth.GetUserRolesResp, error) {

	req := &auth.GetUserRolesReq{
		UserId: userId,
	}

	resp, err := authClient.GetUserRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GrantPermissionToRole 给角色分配权限
func GrantPermissionToRole(ctx context.Context, roleId, permissionId uint64) (bool, error) {
	req := &auth.GrantPermissionToRoleReq{
		RoleId:       roleId,
		PermissionId: permissionId,
	}

	_, err := authClient.GrantPermissionToRole(ctx, req)
	return err == nil, err
}

// RevokePermissionFromRole 移除角色的权限
func RevokePermissionFromRole(ctx context.Context, roleId, permissionId uint64) (bool, error) {
	req := &auth.RevokePermissionFromRoleReq{
		RoleId:       roleId,
		PermissionId: permissionId,
	}
	_, err := authClient.RevokePermissionFromRole(ctx, req)
	return err == nil, err
}

// GetRolePermissions 查询角色所有权限
func GetRolePermissions(ctx context.Context, roleId uint64) (*auth.GetRolePermissionsResp, error) {

	req := &auth.GetRolePermissionsReq{
		RoleId: roleId,
	}

	resp, err := authClient.GetRolePermissions(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserPermissions 查询用户所有权限（含角色下权限）
func GetUserPermissions(ctx context.Context, userId uint64) (*auth.GetUserPermissionsResp, error) {

	req := &auth.GetUserPermissionsReq{
		UserId: userId,
	}

	resp, err := authClient.GetUserPermissions(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// HasPermission 判断用户是否拥有某个权限
func HasPermission(ctx context.Context, userId uint64, permissionCode string) (*auth.HasPermissionResp, error) {

	req := &auth.HasPermissionReq{
		UserId:         userId,
		PermissionCode: permissionCode,
	}

	resp, err := authClient.HasPermission(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CanAccess 判断用户是否可对某资源执行操作
func CanAccess(ctx context.Context, userId uint64, resource, action string) (*auth.CanAccessResp, error) {

	req := &auth.CanAccessReq{
		UserId:   userId,
		Resource: resource,
		Action:   action,
	}

	resp, err := authClient.CanAccess(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListRoles 查询所有角色
func ListRoles(ctx context.Context, page, PageSize uint32) (*auth.ListRolesResp, error) {
	req := &auth.ListRolesReq{
		Page:     page,
		PageSize: PageSize,
	}
	resp, err := authClient.ListRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListPermissions 查询所有权限
func ListPermissions(ctx context.Context, page, pageSize uint32) (*auth.ListPermissionsResp, error) {

	req := &auth.ListPermissionsReq{
		Page:     page,
		PageSize: pageSize,
	}
	resp, err := authClient.ListPermissions(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// IsSuperAdmin 判断是否超级管理员
func IsSuperAdmin(ctx context.Context, userId uint64) (*auth.IsSuperAdminResp, error) {

	req := &auth.IsSuperAdminReq{
		UserId: userId,
	}

	resp, err := authClient.IsSuperAdmin(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
