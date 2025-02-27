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

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct {
    db *gorm.DB
}

// 验证添加商品请求
func validateAddItemReq(req *cart.AddItemReq) error {
    if req == nil {
        return fmt.Errorf("request is nil")
    }
    if req.Item == nil {
        return fmt.Errorf("item is nil")
    }
    if req.Item.Quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }
    if req.UserId == 0 {
        return fmt.Errorf("user id is required")
    }
    return nil
}

// 检查是否为重复键错误
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
            Num:      int(req.Item.Quantity),
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
                    return fmt.Errorf("failed to update cart item")
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
        return nil, fmt.Errorf("user id is required")
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
            ProductId: uint32(item.ProductID),
            Quantity: int32(item.Num),
        }
    }
    return resp, nil
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
    if req.UserId == 0 {
        return nil, fmt.Errorf("user id is required")
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
        return nil, fmt.Errorf("user id is required")
    }
    if req.ProductId == 0 {
        return nil, fmt.Errorf("product id is required")
    }
    if req.Num <= 0 {
        return nil, fmt.Errorf("invalid quantity")
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