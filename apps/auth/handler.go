package main

import (
	"context"
	"encoding/json"
	"github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"math/rand"
	"strconv"
	"time"
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

	//根据userId查询roleIds
	userRolesList, err := s.getUserRolesFromUserId(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	//封装为roleIds
	idList := make([]uint64, 0, len(userRolesList))

	for _, v := range userRolesList {
		idList = append(idList, v.RoleID)
	}

	//根据roleIds获取角色列表
	roleList, err := s.getRolesFromRoleIds(ctx, idList, req.UserId)

	if err != nil {
		return nil, RoleListSearchError
	}

	//组装用户响应
	resList := make([]*auth.Role, 0)

	for _, v := range roleList {
		t := &auth.Role{
			Id:          v.ID,
			RoleName:    v.RoleName,
			Description: v.Description,
			Status:      v.Status == 1,
			CreatedBy:   v.CreatedBy,
		}
		resList = append(resList, t)
	}

	return &auth.GetUserRolesResp{Roles: resList}, nil
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

	//根据roleId查找permIds
	RolePermissionList, err := s.getRolePermsFromRoleId(ctx, req.RoleId)

	if err != nil {
		return nil, SearchRolePermissionError
	}

	//封装为permIds
	idList := make([]uint64, 0, len(RolePermissionList))

	for _, v := range RolePermissionList {
		idList = append(idList, v.PermissionID)
	}

	//根据permIds求得权限列表
	permsList, err := s.getPermsFromPermIds(ctx, idList, req.RoleId)

	if err != nil {
		return nil, err
	}

	//封装权限响应
	pList := make([]*auth.Permission, 0, len(permsList))

	for _, v := range permsList {
		t := &auth.Permission{
			Id:             v.ID,
			PermissionName: v.PermissionName,
			Description:    v.Description,
			ParentId:       v.ParentID,
			Type:           v.Type,
			Resource:       v.Resource,
			Method:         v.Method,
			Status:         v.Status == 1,
			Code:           v.Code,
		}
		pList = append(pList, t)
	}

	return &auth.GetRolePermissionsResp{Permissions: pList}, nil
}

// GetUserPermissions 查询用户所有权限（含角色下权限）
func (s *AuthServiceImpl) GetUserPermissions(ctx context.Context, req *auth.GetUserPermissionsReq) (resp *auth.GetUserPermissionsResp, err error) {

	userRoleList, err := s.getUserRolesFromUserId(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	roleIds := make([]uint64, 0, len(userRoleList))

	for _, v := range userRoleList {
		roleIds = append(roleIds, v.RoleID)
	}

	permsMap := make(map[uint64]*auth.Permission)

	for _, v := range roleIds {

		permList, err := s.GetRolePermissions(ctx, &auth.GetRolePermissionsReq{RoleId: v})

		if err != nil {
			return nil, err
		}

		for _, perm := range permList.Permissions {
			permsMap[perm.Id] = perm
		}

	}

	permsList := make([]*auth.Permission, 0, len(permsMap))

	for _, v := range permsMap {
		t := &auth.Permission{
			Id:             v.Id,
			PermissionName: v.PermissionName,
			Description:    v.Description,
			ParentId:       v.ParentId,
			Type:           v.Type,
			Resource:       v.Resource,
			Method:         v.Method,
			Status:         v.Status,
			Code:           v.Code,
		}
		permsList = append(permsList, t)
	}

	return &auth.GetUserPermissionsResp{Permissions: permsList}, nil
}

// HasPermission 判断用户是否拥有某个权限
func (s *AuthServiceImpl) HasPermission(ctx context.Context, req *auth.HasPermissionReq) (resp *auth.HasPermissionResp, err error) {
	// TODO: Your code here...
	key := util.TakeKey(permission, "is_allowed", req.UserId, req.PermissionCode)

	resInLocal, has := LocalCache.Get(key)

	if has {
		if cacheV1, ok := resInLocal.(bool); ok {
			return &auth.HasPermissionResp{Ok: cacheV1}, nil
		} else {
			LocalCache.Del(key)
		}
	}

	resInRedis, err := Rds.Get(ctx, key).Result()

	if err == nil {
		if parseBool, err := strconv.ParseBool(resInRedis); err == nil {
			LocalCache.Set(key, parseBool, 1)
			return &auth.HasPermissionResp{Ok: parseBool}, nil
		} else {
			LocalCache.Del(key)
		}
	}

	return
}

// CanAccess 判断用户是否可对某资源执行操作
func (s *AuthServiceImpl) CanAccess(ctx context.Context, req *auth.CanAccessReq) (resp *auth.CanAccessResp, err error) {
	// TODO: Your code here...

	return
}

// ListRoles 查询所有角色
func (s *AuthServiceImpl) ListRoles(ctx context.Context, req *auth.ListRolesReq) (resp *auth.ListRolesResp, err error) {

	if req.PageSize < 1 || req.Page < 1 {
		return nil, BadPageOrPageSize
	}

	offset := (req.Page - 1) * req.PageSize

	listQuery := util.ListCacheComponent[uint64, models.Role]{
		Rds:             Rds,
		Ctx:             ctx,
		IdListKey:       util.TakeKey(role, "list", util.Md5Hash(util.TakeKey(req.Page, req.PageSize))),
		DetailKeyPrefix: util.TakeKey(role),
		FuncName:        "ListRoles",
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FullQueryExec: func() ([]models.Role, error) {
			list := make([]models.Role, 0)
			if err := DB.
				Model(&models.Role{}).
				Offset(int(offset)).
				Limit(int(req.PageSize)).
				Find(&list).Error; err != nil {
				return nil, err
			}
			return list, nil
		},
		Expires:     time.Duration(5+rand.Intn(3)) * time.Minute,
		MaxLostRate: 30,
		Sort:        nil,
		IdName:      "id",
		DB:          DB,
	}

	cache, err := listQuery.QueryListWithCache()

	if err != nil {
		return nil, err
	}

	roles := make([]*auth.Role, 0, len(cache))

	for _, v := range cache {
		t := &auth.Role{
			Id:          v.ID,
			RoleName:    v.RoleName,
			Description: v.Description,
			Status:      v.Status == 1,
			CreatedBy:   v.CreatedBy,
			Code:        v.Code,
		}
		roles = append(roles, t)
	}

	return &auth.ListRolesResp{Roles: roles}, nil
}

// ListPermissions 查询所有权限
func (s *AuthServiceImpl) ListPermissions(ctx context.Context, req *auth.ListPermissionsReq) (resp *auth.ListPermissionsResp, err error) {

	if req.PageSize < 1 || req.Page < 1 {
		return nil, BadPageOrPageSize
	}

	offset := (req.Page - 1) * req.PageSize

	listQuery := util.ListCacheComponent[uint64, models.Permission]{
		Rds:             Rds,
		Ctx:             ctx,
		IdListKey:       util.TakeKey(permission, "list", util.Md5Hash(util.TakeKey(req.Page, req.PageSize))),
		DetailKeyPrefix: util.TakeKey(permission),
		FuncName:        "ListPermissions",
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FullQueryExec: func() ([]models.Permission, error) {
			list := make([]models.Permission, 0)
			if err := DB.
				Model(&models.Permission{}).
				Offset(int(offset)).
				Limit(int(req.PageSize)).
				Find(&list).Error; err != nil {
				return nil, err
			}
			return list, nil
		},
		Expires:     time.Duration(5+rand.Intn(3)) * time.Minute,
		MaxLostRate: 30,
		Sort:        nil,
		IdName:      "id",
		DB:          DB,
	}

	cache, err := listQuery.QueryListWithCache()

	if err != nil {
		return nil, err
	}

	perms := make([]*auth.Permission, 0, len(cache))

	for _, v := range cache {
		t := &auth.Permission{
			Id:             v.ID,
			PermissionName: v.PermissionName,
			Description:    v.Description,
			ParentId:       v.ParentID,
			Type:           v.Type,
			Resource:       v.Resource,
			Method:         v.Method,
			Status:         v.Status == 1,
			Code:           v.Code,
		}
		perms = append(perms, t)
	}

	return &auth.ListPermissionsResp{Perms: perms}, nil
}

// IsSuperAdmin 判断是否超级管理员
// 这种允许为用户表添加字段的行为适用于电商系统这种读多写少的模式,在其他系统中慎用
func (s *AuthServiceImpl) IsSuperAdmin(ctx context.Context, req *auth.IsSuperAdminReq) (resp *auth.IsSuperAdminResp, err error) {

	key := util.TakeKey("is_super_admin", req.UserId)

	//先查询本地缓存
	res, isExist := LocalCache.Get(key)

	if isExist {
		if b, ok := res.(bool); ok {
			return &auth.IsSuperAdminResp{
				Ok: b,
			}, nil
		} else {
			LocalCache.Del(key)
		}
	}

	//查询redis缓存
	result, err := Rds.Get(ctx, key).Result()

	if err == nil {
		parseBool, Err := strconv.ParseBool(result)

		if Err != nil {
			Rds.Del(ctx, key)
		} else {

			LocalCache.Set(key, parseBool, 1)

			return &auth.IsSuperAdminResp{
				Ok: parseBool,
			}, nil
		}
	}

	//查询数据库

	var u models.User

	if err = DB.Where(&models.User{}).Where("id = ?", req.UserId).First(&u).Error; err != nil {
		return nil, SearchSuperError
	}

	//设置缓存
	s.setDoubleCaches(ctx, key, u.IsSuperAdmin, time.Duration(10+rand.Intn(3))*time.Minute)

	return &auth.IsSuperAdminResp{Ok: u.IsSuperAdmin}, nil
}
