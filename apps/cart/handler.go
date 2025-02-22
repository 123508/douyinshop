package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cart "github.com/123508/douyinshop/kitex_gen/cart"
	"github.com/123508/douyinshop/pkg/models"
	"gorm.io/gorm"
)

type CartServiceImpl struct {
	db *gorm.DB
}

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}

func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error) {
	if req.Item == nil || req.Item.Quantity <= 0 {
		return nil, fmt.Errorf("invalid item or quantity")
	}

	userID := req.UserId
	productID := req.Item.ProductId
	quantity := int(req.Item.Quantity)

	newCart := models.Cart{
		UserId:    userID,
		ProductID: productID,
		Num:       quantity,
	}

	// 尝试创建，若唯一约束冲突则更新数量
	err := s.db.Create(&newCart).Error
	if err != nil {
		if IsDuplicateKeyError(err) {
			// 更新现有记录的数量
			result := s.db.Model(&models.Cart{}).
				Where("user_id = ? AND product_id = ?", userID, productID).
				Update("num", gorm.Expr("num + ?", quantity))
			if result.Error != nil {
				return nil, result.Error
			}
		} else {
			return nil, err
		}
	}
	return &cart.AddItemResp{}, nil
}

func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error) {
	var cartItems []models.Cart
	result := s.db.Where("user_id = ?", req.UserId).Find(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}

	pbItems := make([]*cart.CartItem, len(cartItems))
	for i, item := range cartItems {
		pbItems[i] = &cart.CartItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Num),
		}
	}

	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.UserId,
			Items:  pbItems,
		},
	}, nil
}

func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	result := s.db.Where("user_id = ?", req.UserId).Delete(&models.Cart{})
	if result.Error != nil {
		return nil, result.Error
	}
	return &cart.EmptyCartResp{}, nil
}

func (s *CartServiceImpl) DeleteItem(ctx context.Context, req *cart.DeleteItemReq) (*cart.EmptyCartResp, error) {
	if req.Num <= 0 {
		return nil, fmt.Errorf("invalid num")
	}

	var cartItem models.Cart
	result := s.db.Where("user_id = ? AND product_id = ?", req.UserId, req.ProductId).First(&cartItem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &cart.EmptyCartResp{}, nil
		}
		return nil, result.Error
	}

	newNum := cartItem.Num - int(req.Num)
	if newNum <= 0 {
		if err := s.db.Delete(&cartItem).Error; err != nil {
			return nil, err
		}
	} else {
		cartItem.Num = newNum
		if err := s.db.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	}
	return &cart.EmptyCartResp{}, nil
}
