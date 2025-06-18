package main

import (
	"context"
	"encoding/json"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

const (
	serviceName = "auth"
	permission  = "permission"
	role        = "role"
	user        = "user"
)

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		util.LogError("打开MySQL连接失败", "connectWithMySQL", err)
		panic(err)
	}
	return DB
}

var Rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", err)
		panic(err)
	}
	return rds
}

var LocalCache = connectWithLocalCache()

func connectWithLocalCache() *util.RistrettoCache {
	cache, err := util.NewRistrettoCache()
	if err != nil {
		util.LogError("本地缓存开启失败", "connectWithLocalCache", err)
		panic(err)
	}
	return cache
}

// 辅助函数：设置两级缓存
func (s *AuthServiceImpl) setDoubleCaches(ctx context.Context, key string, value interface{}, randTTL time.Duration) {
	// Redis缓存：使用随机TTL防止雪崩
	Rds.Set(ctx, key, value, randTTL)

	// 本地缓存：设置比Redis更短的TTL
	LocalCache.Set(key, value, 1)
}

func (s *AuthServiceImpl) cleanSuperAdminCache(ctx context.Context, userId uint64) {
	key := util.TakeKey("is_super_admin", userId)
	Rds.Del(ctx, key)
	LocalCache.Del(key)
}

// 查询用户对应的角色Id列表
func (s *AuthServiceImpl) getUserRolesFromUserId(ctx context.Context, userId uint64) ([]models.UserRole, error) {

	userRoleListQuery := util.SimpleCacheComponent[uint64, []models.UserRole]{
		Rds:       Rds,
		Ctx:       ctx,
		Key:       util.TakeKey(user, role, userId),
		FuncName:  "getUserRolesFromUserId",
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		QueryExec: func() ([]models.UserRole, error) {

			UserRoleList := make([]models.UserRole, 0)

			if err := DB.Model(&models.UserRole{}).
				Where("user_id = ?", userId).
				Find(&UserRoleList).Error; err != nil {
				return nil, err
			}
			return UserRoleList, nil
		},
		Expires: time.Duration(5+rand.Intn(3)) * time.Minute,
	}

	userRolesList, err := userRoleListQuery.QueryWithCache()

	if err != nil {
		util.LogError("用户角色查询错误", "getUserRolesFromUserId", err)
		return nil, UserRoleSearchError
	}

	return userRolesList, nil
}

// 根据查询的roleId列表查询roles
func (s *AuthServiceImpl) getRolesFromRoleIds(ctx context.Context, roleIds []uint64, userId uint64) ([]models.Role, error) {
	//查询角色列表
	rolesListQuery := util.ListCacheComponent[uint64, models.Role]{
		Rds:             Rds,
		Ctx:             ctx,
		IdListKey:       util.TakeKey(serviceName, role, util.Md5Hash(util.TakeKey(userId))),
		DetailKeyPrefix: util.TakeKey(role),
		FuncName:        "getRolesFromRoleIds",
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FullQueryExec: func() ([]models.Role, error) {

			roles := make([]models.Role, 0)

			if err := DB.Model(&models.Role{}).Where("id in ?", roleIds).
				Find(&roles).Error; err != nil {
				return nil, err
			}
			return roles, nil
		},
		Expires:     time.Duration(5+rand.Intn(3)) * time.Minute,
		MaxLostRate: 30,
		Sort:        nil,
		IdName:      "id",
		DB:          DB,
	}

	roleList, err := rolesListQuery.QueryListWithCache()

	if err != nil {
		util.LogError("查询角色列表错误", "getRolesFromRoleIds", err)
		return nil, RoleListSearchError
	}

	return roleList, nil
}

// 查询角色对应的权限Id列表
func (s *AuthServiceImpl) getRolePermsFromRoleId(ctx context.Context, roleId uint64) ([]models.RolePermission, error) {

	userRoleListQuery := util.SimpleCacheComponent[uint64, []models.RolePermission]{
		Rds:       Rds,
		Ctx:       ctx,
		Key:       util.TakeKey(role, permission, roleId),
		FuncName:  "getRolePermsFromRoleId",
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		QueryExec: func() ([]models.RolePermission, error) {

			RolePermissionList := make([]models.RolePermission, 0)

			if err := DB.Model(&models.RolePermission{}).
				Where("role_id = ?", roleId).
				Find(&RolePermissionList).Error; err != nil {
				return nil, err
			}
			return RolePermissionList, nil
		},
		Expires: time.Duration(5+rand.Intn(3)) * time.Minute,
	}

	RolePermissionList, err := userRoleListQuery.QueryWithCache()

	if err != nil {
		util.LogError("角色权限查询错误", "getRolePermsFromRoleId", err)
		return nil, SearchRolePermissionError
	}
	return RolePermissionList, nil
}

func (s *AuthServiceImpl) getPermsFromPermIds(ctx context.Context, permIds []uint64, roleId uint64) ([]models.Permission, error) {
	//查询权限列表
	rolesListQuery := util.ListCacheComponent[uint64, models.Permission]{
		Rds:             Rds,
		Ctx:             ctx,
		IdListKey:       util.TakeKey(serviceName, permission, util.Md5Hash(util.TakeKey(roleId))),
		DetailKeyPrefix: util.TakeKey(permission),
		FuncName:        "getPermsFromPermIds",
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FullQueryExec: func() ([]models.Permission, error) {
			perms := make([]models.Permission, 0)

			if err := DB.Model(&models.Permission{}).
				Where("id in ?", permIds).
				Find(&perms).Error; err != nil {
				return nil, err
			}
			return perms, nil
		},
		Expires:     time.Duration(5+rand.Intn(3)) * time.Minute,
		MaxLostRate: 30,
		Sort:        nil,
		IdName:      "id",
		DB:          DB,
	}

	permsList, err := rolesListQuery.QueryListWithCache()

	if err != nil {
		util.LogError("查询角色权限错误", "getPermsFromPermIds", err)
		return nil, SearchRolePermissionError
	}
	return permsList, nil
}
