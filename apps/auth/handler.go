package main

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// AssignRoleToUser 给用户分配角色
func (s *AuthServiceImpl) AssignRoleToUser(ctx context.Context, req *auth.AssignRoleToUserReq) (resp *auth.Empty, err error) {
	// TODO: Your code here...
	return
}

// RemoveRoleFromUser 移除用户的角色
func (s *AuthServiceImpl) RemoveRoleFromUser(ctx context.Context, req *auth.RemoveRoleFromUserReq) (resp *auth.Empty, err error) {
	// TODO: Your code here...
	return
}

// GetUserRoles 查询用户所有角色
func (s *AuthServiceImpl) GetUserRoles(ctx context.Context, req *auth.GetUserRolesReq) (resp *auth.GetUserRolesResp, err error) {
	// TODO: Your code here...
	return
}

// GrantPermissionToRole 给角色分配权限
func (s *AuthServiceImpl) GrantPermissionToRole(ctx context.Context, req *auth.GrantPermissionToRoleReq) (resp *auth.Empty, err error) {
	// TODO: Your code here...
	return
}

// RevokePermissionFromRole 移除角色的权限
func (s *AuthServiceImpl) RevokePermissionFromRole(ctx context.Context, req *auth.RevokePermissionFromRoleReq) (resp *auth.Empty, err error) {
	// TODO: Your code here...
	return
}

// GetRolePermissions 查询角色所有权限
func (s *AuthServiceImpl) GetRolePermissions(ctx context.Context, req *auth.GetRolePermissionsReq) (resp *auth.GetRolePermissionsResp, err error) {
	// TODO: Your code here...
	return
}

// GetUserPermissions 查询用户所有权限（含角色下权限）
func (s *AuthServiceImpl) GetUserPermissions(ctx context.Context, req *auth.GetUserPermissionsReq) (resp *auth.GetUserPermissionsResp, err error) {
	// TODO: Your code here...
	return
}

// HasPermission 判断用户是否拥有某个权限
func (s *AuthServiceImpl) HasPermission(ctx context.Context, req *auth.HasPermissionReq) (resp *auth.HasPermissionResp, err error) {
	// TODO: Your code here...
	return
}

// CanAccess 判断用户是否可对某资源执行操作
func (s *AuthServiceImpl) CanAccess(ctx context.Context, req *auth.CanAccessReq) (resp *auth.CanAccessResp, err error) {
	// TODO: Your code here...
	return
}

// ListRoles 查询所有角色
func (s *AuthServiceImpl) ListRoles(ctx context.Context, req *auth.Empty) (resp *auth.ListRolesResp, err error) {
	// TODO: Your code here...
	return
}

// ListPermissions 查询所有权限
func (s *AuthServiceImpl) ListPermissions(ctx context.Context, req *auth.Empty) (resp *auth.ListPermissionsResp, err error) {
	// TODO: Your code here...
	return
}

// IsSuperAdmin 判断是否超级管理员
func (s *AuthServiceImpl) IsSuperAdmin(ctx context.Context, req *auth.IsSuperAdminReq) (resp *auth.IsSuperAdminResp, err error) {
	// TODO: Your code here...
	return
}
