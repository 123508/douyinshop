package repository

import (
	"context"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	GetDB() *gorm.DB
	CreateUser(ctx context.Context, user *models.User, uLogin *models.UserLogin) (uint64, error)
	CompareUserPassword(ctx context.Context, userId uint64, password string) (bool, error)
	GetUserInfoFromId(ctx context.Context, userId uint64) (*models.User, error)
	GetUserInfoFromEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userId uint64) error
	GetUserList(ctx context.Context, page, pageSize int, filterSql string, filterParams []interface{}) ([]models.User, error)
	UpdateUserPassword(ctx context.Context, userId uint64, password string) error
	RemoveEmail(ctx context.Context, userId uint64) error
	SetEmail(ctx context.Context, userId uint64, email string) error
	FreezeUser(ctx context.Context, userId uint64) error
	UnfreezeUser(ctx context.Context, userId uint64) error
}

type RepoImpl struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) UserRepository {
	return &RepoImpl{
		DB: db,
	}
}

func (r *RepoImpl) GetDB() *gorm.DB {
	return r.DB
}

func (r *RepoImpl) CreateUser(ctx context.Context, user *models.User, uLogin *models.UserLogin) (uint64, error) {

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		//创建用户对象
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		uLogin.UserId = user.ID

		if err := tx.Create(uLogin).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		util.LogError("创建用户对象错误", "CreateUser", err)
		return 0, err
	}

	return user.ID, nil
}

func (r *RepoImpl) CompareUserPassword(ctx context.Context, userId uint64, password string) (bool, error) {
	var res models.UserLogin
	if err := r.DB.Model(&models.UserLogin{}).
		WithContext(ctx).
		Where("user_id = ? and password = ?", userId, password).
		Find(&res).Error; err != nil {
		util.LogError("查询用户密码错误", "GetUserPassword", err)
		return false, err
	}

	return res.ID != 0, nil
}

func (r *RepoImpl) GetUserInfoFromId(ctx context.Context, userId uint64) (*models.User, error) {
	var row models.User
	if err := r.DB.Model(&models.User{}).WithContext(ctx).Where("id = ?", userId).First(&row).Error; err != nil {
		util.LogError("通过id获取用户信息", "GetUserInfoFromId", err)
		return nil, err
	}
	return &row, nil
}

func (r *RepoImpl) GetUserInfoFromEmail(ctx context.Context, email string) (*models.User, error) {
	//从users表中读取Email信息
	var row models.User
	if err := r.DB.Model(&models.User{}).WithContext(ctx).Where("email = ?", email).Find(&row).Error; err != nil {
		util.LogError("通过email获取用户信息", "GetUserInfoFromEmail", err)
		return nil, err
	}
	return &row, nil
}

func (r *RepoImpl) UpdateUser(ctx context.Context, u *models.User) error {

	updates := make(map[string]interface{})

	updates["gender"] = u.Gender

	if u.Phone != "" {
		updates["phone"] = u.Phone
	}

	if u.Name != "" {
		updates["name"] = u.Name
	}

	if u.Avatar != "" {
		updates["avatar"] = u.Avatar
	}

	if err := r.DB.WithContext(ctx).Model(&models.User{}).
		Where("id=?", u.ID).
		Updates(updates).Error; err != nil {
		util.LogError("更新用户出错", "UpdateUser", err)
		return err
	}
	return nil
}

func (r *RepoImpl) DeleteUser(ctx context.Context, UserId uint64) error {

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).
			Where("id = ?", UserId).
			Update("phone", nil).
			Update("email", nil).
			Delete(&models.User{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.UserLogin{}).
			Where("user_id=?", UserId).
			Delete(&models.UserLogin{UserId: UserId}).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		util.LogError("删除用户异常", "DeleteUser", err)
		return err
	}
	return nil
}

func (r *RepoImpl) GetUserList(ctx context.Context,
	page, pageSize int, filterSql string, filterParams []interface{}) ([]models.User, error) {

	var users []models.User
	offset := (page - 1) * pageSize
	if err := r.DB.WithContext(ctx).Model(&models.User{}).
		Where(filterSql, filterParams...).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		util.LogError("查询用户列表失败", "GetUserList", err)
		return nil, err
	}
	return users, nil
}

func (r *RepoImpl) UpdateUserPassword(ctx context.Context, userId uint64, password string) error {
	where := r.DB.WithContext(ctx).Model(&models.UserLogin{}).Where("user_id = ?", userId)

	if err := where.Update("password", password).Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("更新用户密码错误", "UpdateUserPassword", err)
		return err
	}

	return nil
}

func (r *RepoImpl) RemoveEmail(ctx context.Context, userId uint64) error {
	where := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userId)

	if err := where.Update("email", "").Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("解绑邮箱失败", "RemoveEmail", err)
		return err
	}
	return nil
}

func (r *RepoImpl) SetEmail(ctx context.Context, userId uint64, email string) error {
	where := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userId)

	if err := where.Update("email", email).Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("解绑邮箱失败", "RemoveEmail", err)
		return err
	}
	return nil
}

func (r *RepoImpl) FreezeUser(ctx context.Context, userId uint64) error {

	if err := r.DB.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userId).
		Update("status", 1).
		Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("冻结用户失败", "FreezeUser", err)
		return err
	}

	return nil
}

func (r *RepoImpl) UnfreezeUser(ctx context.Context, userId uint64) error {
	if err := r.DB.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userId).
		Update("status", 0).
		Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("解冻用户失败", "UnfreezeUser", err)
		return err
	}

	return nil
}
