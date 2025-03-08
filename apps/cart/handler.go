package main

import (
	"context"
	"errors"
	"github.com/123508/douyinshop/pkg/errorno"
	"strings"

	cart "github.com/123508/douyinshop/kitex_gen/cart"
	"github.com/123508/douyinshop/pkg/models"
	"gorm.io/gorm"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct {
	db *gorm.DB
}

var NilRequestError = &errorno.BasicMessageError{Code: 400, Message: "请求为空"}

var NilItemError = &errorno.BasicMessageError{Code: 400, Message: "商品为空"}

var NegativeQuantityError = &errorno.BasicMessageError{Code: 400, Message: "数量必须为正数"}

var NilUserIdError = &errorno.BasicMessageError{Code: 400, Message: "用户id不能为空"}

var FailToUpdateCart = &errorno.BasicMessageError{Code: 403, Message: "无法更新购物车"}

var NilProductIdError = &errorno.BasicMessageError{Code: 400, Message: "商品id不能为空"}

// 验证添加商品请求
func validateAddItemReq(req *cart.AddItemReq) error {
	if req == nil {
		return NilRequestError
	}
	if req.Item == nil {
		return NilItemError
	}
	if req.Item.Quantity <= 0 {
		return NegativeQuantityError
	}
	if req.UserId == 0 {
		return NilUserIdError
	}
	return nil
}

// IsDuplicateKeyError 检查是否为重复键错误
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error) {
	// 1. 输入验证
	if err := validateAddItemReq(req); err != nil {
		return nil, err
	}

	// 2. 使用事务确保原子性
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newCart := models.Cart{
			UserId:    req.UserId,
			ProductID: req.Item.ProductId,
			Num:       int(req.Item.Quantity),
		}

		// 3. 尝试创建
		err := tx.Create(&newCart).Error
		if err != nil {
			if IsDuplicateKeyError(err) {
				// 4. 更新现有记录
				result := tx.Model(&models.Cart{}).
					Where("user_id = ? AND product_id = ?", req.UserId, req.Item.ProductId).
					Update("num", gorm.Expr("num + ?", req.Item.Quantity))
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return FailToUpdateCart
				}
				return nil
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &cart.AddItemResp{}, nil
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error) {
	if req.UserId == 0 {
		return nil, NilUserIdError
	}

	resp := &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.UserId,
			Items:  make([]*cart.CartItem, 0),
		},
	}

	// 查询用户的所有购物车商品
	var cartItems []models.Cart
	if err := s.db.WithContext(ctx).Where("user_id = ?", req.UserId).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	// 预分配切片容量，提高性能
	resp.Cart.Items = make([]*cart.CartItem, len(cartItems))
	for i, item := range cartItems {
		resp.Cart.Items[i] = &cart.CartItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Num),
		}
	}
	return resp, nil
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	if req.UserId == 0 {
		return nil, NilUserIdError
	}

	result := s.db.WithContext(ctx).Where("user_id = ?", req.UserId).Delete(&models.Cart{})
	if result.Error != nil {
		return nil, result.Error
	}

	return &cart.EmptyCartResp{}, nil
}

// DeleteItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) DeleteItem(ctx context.Context, req *cart.DeleteItemReq) (*cart.EmptyCartResp, error) {
	if req.UserId == 0 {
		return nil, NilUserIdError
	}
	if req.ProductId == 0 {
		return nil, NilProductIdError
	}
	if req.Num <= 0 {
		return nil, NegativeQuantityError
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cartItem models.Cart
		result := tx.Where("user_id = ? AND product_id = ?", req.UserId, req.ProductId).First(&cartItem)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil // 商品不存在，视为成功
			}
			return result.Error
		}

		if uint32(cartItem.Num) <= req.Num {
			// 如果删除数量大于等于现有数量，删除整个记录
			if err := tx.Delete(&cartItem).Error; err != nil {
				return err
			}
		} else {
			// 否则减少数量
			cartItem.Num -= int(req.Num)
			if err := tx.Save(&cartItem).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &cart.EmptyCartResp{}, nil
}
