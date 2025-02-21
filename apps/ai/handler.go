package main

import (
    "context"
    "encoding/json"
    "time" 
    "strings"  
    ai "github.com/123508/douyinshop/kitex_gen/ai"
    "github.com/123508/douyinshop/kitex_gen/order/orderservice"
    aiutil "github.com/123508/douyinshop/pkg/ai"
    "github.com/123508/douyinshop/pkg/errors"
    "github.com/cloudwego/kitex/pkg/klog" 
    "github.com/123508/douyinshop/pkg/config"
)

const (
    maxRequestLength = 1000
    defaultTimeout  = 10 * time.Second
)

// AiServiceImpl implements the last service interface defined in the IDL.
type AiServiceImpl struct{
    douyinAI     *aiutil.DouyinAI
    orderClient  orderservice.Client
    timeout      time.Duration
}

// NewAiServiceImpl 创建服务实现实例
func NewAiServiceImpl(orderClient orderservice.Client) *AiServiceImpl {
    timeout := config.Conf.AiConfig.Timeout
    if timeout <= 0 {
        timeout = defaultTimeout
    }
    
    return &AiServiceImpl{
        douyinAI: aiutil.NewDouyinAI(),
        orderClient: orderClient,
        timeout: timeout,
    }
}

// Close 关闭资源
func (s *AiServiceImpl) Close() error {
    if s.douyinAI != nil {
        return s.douyinAI.Close()
    }
    return nil
}

// OrderQuery implements the AiServiceImpl interface.
func (s *AiServiceImpl) OrderQuery(ctx context.Context, req *ai.OrderQueryReq) (resp *ai.OrderQueryResp, err error) {
    // 参数验证
    if req == nil {
        return nil, &errors.BasicMessageError{Message: "请求不能为空"}
    }
    if req.OrderId == "" {
        return nil, &errors.BasicMessageError{Message: "订单ID不能为空"}
    }

    // 设置超时
    ctx, cancel := context.WithTimeout(ctx, s.timeout)
    defer cancel()

    // 调用订单服务
    orderResp, err := s.orderClient.GetOrder(ctx, &orderservice.GetOrderReq{OrderId: req.OrderId})  
    if err != nil {
        klog.Errorf("GetOrder failed for orderID %s: %v", req.OrderId, err)
        return nil, errors.WrapWithMessage(err, "获取订单信息失败")
    }
    if orderResp == nil || orderResp.Order == nil {
        klog.Errorf("GetOrder returned nil for orderID %s", req.OrderId)
        return nil, errors.New("订单信息不存在")
    }

    // 记录操作日志
    klog.Infof("Processing order query for orderID: %s", req.OrderId)

    // 手动构建orderMap避免序列化开销
    orderMap := map[string]interface{}{
        "id":     orderResp.Order.Id,
        "status": orderResp.Order.Status,
        "amount": orderResp.Order.TotalAmount,
        "create_time":  orderResp.Order.CreateTime,
        "items":  orderResp.Order.Items,
        "user_id": orderResp.Order.UserId,
    }

    // 调用AI格式化
    response, err := s.douyinAI.FormatOrderDetails(orderMap)
    if err != nil {
        klog.Errorf("FormatOrderDetails failed for orderID %s: %v", req.OrderId, err)
        return nil, errors.Wrap(err, "格式化订单详情失败")
    }

    klog.Infof("Successfully processed order query for orderID: %s", req.OrderId)
    return &ai.OrderQueryResp{Response: response}, nil
}

func (s *AiServiceImpl) AutoPlaceOrder(ctx context.Context, req *ai.AutoPlaceOrderReq) (resp *ai.AutoPlaceOrderResp, err error) {
    // 参数验证
    if req == nil {
        return nil, &errors.BasicMessageError{Message: "请求不能为空"}
    }
    if req.UserId <= 0 {
        return nil, &errors.BasicMessageError{Message: "无效的用户ID"}
    }
    if request := strings.TrimSpace(req.Request); request == "" {
        return nil, &errors.BasicMessageError{Message: "请求内容不能为空"}
    } else if len(request) > maxRequestLength {
        return nil, &errors.BasicMessageError{Message: "请求内容超出长度限制"}
    }

    klog.Infof("Processing auto place order request for userID: %d", req.UserId)

    // AI分析请求
    recommendations, err := s.douyinAI.AnalyzeOrderRequest(req.Request)
    if err != nil {
        klog.Errorf("AnalyzeOrderRequest failed for userID %d: %v", req.UserId, err)
        return nil, errors.Wrap(err, "分析订单请求失败")
    }
    if len(recommendations) == 0 {
        return nil, errors.New("未找到推荐商品")
    }

    // 构建订单项
    items := make([]*orderservice.OrderItem, 0, len(recommendations))
    for _, rec := range recommendations {
        if err := validateRecommendation(rec); err != nil {
            klog.Warnf("Invalid recommendation for userID %d: %+v, error: %v", req.UserId, rec, err)
            continue
        }
        items = append(items, &orderservice.OrderItem{
            ProductId: rec.ProductID,
            Quantity:  int32(rec.Quantity),
        })
    }
    if len(items) == 0 {
        return nil, errors.New("没有有效的订单项")
    }

    // 创建订单
    ctx, cancel := context.WithTimeout(ctx, s.timeout)
    defer cancel()
    
    orderResp, err := s.orderClient.CreateOrder(ctx, &orderservice.CreateOrderReq{
        UserId: req.UserId,
        Items:  items,
    })
    if err != nil {
        klog.Errorf("CreateOrder failed for userID %d: %v", req.UserId, err)
        return nil, errors.Wrap(err, "创建订单失败")
    }
    if orderResp == nil {
        return nil, errors.New("创建订单失败：服务返回为空")
    }

    klog.Infof("Successfully created order %s for userID: %d", orderResp.OrderId, req.UserId)
    return &ai.AutoPlaceOrderResp{OrderId: orderResp.OrderId}, nil
}

func validateRecommendation(rec *aiutil.Recommendation) error {
    if rec == nil {
        return errors.New("推荐信息为空")
    }
    if rec.ProductID <= 0 {
        return errors.New("无效的商品ID")
    }
    if rec.Quantity <= 0 || rec.Quantity > 999 {
        return errors.New("无效的商品数量")
    }
    if rec.Price < 0 {
        return errors.New("无效的商品价格")
    }
    return nil
}