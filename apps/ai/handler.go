package main

import (
    "context"
    "fmt"
    "errors"
    "time" 
    "strings"  
    ai "github.com/123508/douyinshop/kitex_gen/ai"
    "github.com/123508/douyinshop/kitex_gen/order/order_common"
    "github.com/123508/douyinshop/kitex_gen/order/userOrder"
    "github.com/123508/douyinshop/kitex_gen/order/userOrder/orderuserservice"
    aiutil "github.com/123508/douyinshop/pkg/ai"
    "github.com/cloudwego/kitex/pkg/klog" 
    "github.com/123508/douyinshop/pkg/config"
    "strconv"
)
const (
    maxRequestLength = 1000
    defaultTimeout  = 10 * time.Second
)

// AiServiceImpl implements the last service interface defined in the IDL.
type AiServiceImpl struct{
    douyinAI     *aiutil.DouyinAI
    orderClient  orderuserservice.Client
    timeout      time.Duration
}

// NewAiServiceImpl 创建服务实现实例
func NewAiServiceImpl(orderClient orderuserservice.Client) *AiServiceImpl {
    timeout := time.Duration(config.Conf.VolcengineConfig.Timeout) * time.Second
    if timeout <= 0 {
        timeout = defaultTimeout
    }
    
    douyinAI, err := aiutil.NewDouyinAI()
    if err != nil {
        klog.Errorf("Failed to initialize DouyinAI: %v", err)
        return nil, errors.Wrap(err, "初始化AI服务失败")
    }
    
    return &AiServiceImpl{
        douyinAI: douyinAI,
        orderClient: orderClient,
        timeout: timeout,
    }, nil
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
        return nil, errors.New("请求不能为空")
    }
    if req.OrderId == "" {
        return nil, errors.New("订单ID不能为空")
    }

    // 设置超时
    ctx, cancel := context.WithTimeout(ctx, s.timeout)
    defer cancel()

    // 将字符串订单ID转换为uint32
    orderID, err := strconv.ParseUint(req.OrderId, 10, 32)
    if err != nil {
        klog.Errorf("Invalid order ID format %s: %v", req.OrderId, err)
        return nil, errors.New("订单ID格式不正确")
    }

    // 调用订单服务的Detail方法
    orderResp, err := s.orderClient.Detail(ctx, &order_common.OrderReq{OrderId: uint32(orderID)})  
    if err != nil {
        klog.Errorf("Detail failed for orderID %s: %v", req.OrderId, err)
        return nil, fmt.Errorf("获取订单信息失败: %w", err)
    }
    if orderResp == nil || orderResp.Order == nil {
        klog.Errorf("Detail returned nil for orderID %s", req.OrderId)
        return nil, errors.New("订单信息不存在")
    }

    // 记录操作日志
    klog.Infof("Processing order query for orderID: %s", req.OrderId)

    // 构建orderMap
    orderMap := map[string]interface{}{
        "number":     orderResp.Order.Number,
        "status":     orderResp.Order.Status,
        "amount":     orderResp.Order.Amount,
        "user_id":    orderResp.Order.UserId,
        "consignee":  orderResp.Order.Consignee,
        "address":    orderResp.Order.Address,
        "phone":      orderResp.Order.Phone,
        "remark":     orderResp.Order.Remark,
    }

    // 调用AI格式化
    response, err := s.douyinAI.FormatOrderDetails(orderMap)
    if err != nil {
        klog.Errorf("FormatOrderDetails failed for orderID %s: %v", req.OrderId, err)
        return nil, fmt.Errorf("格式化订单详情失败: %w", err)
    }

    klog.Infof("Successfully processed order query for orderID: %s", req.OrderId)
    return &ai.OrderQueryResp{Response: response}, nil
}

func (s *AiServiceImpl) AutoPlaceOrder(ctx context.Context, req *ai.AutoPlaceOrderReq) (resp *ai.AutoPlaceOrderResp, err error) {
    // 参数验证
    if req == nil {
        return nil, errors.New("请求不能为空")
    }
    if req.UserId <= 0 {
        return nil, errors.New("无效的用户ID")
    }
    if request := strings.TrimSpace(req.Request); request == "" {
        return nil, errors.New("请求内容不能为空")
    } else if len(request) > maxRequestLength {
        return nil, errors.New("请求内容超出长度限制")
    }

    klog.Infof("Processing auto place order request for userID: %d", req.UserId)

    // AI分析请求
    recommendations, err := s.douyinAI.AnalyzeOrderRequest(req.Request)
    if err != nil {
        klog.Errorf("AnalyzeOrderRequest failed for userID %d: %v", req.UserId, err)
        return nil, fmt.Errorf("分析订单请求失败: %w", err)
    }
    if len(recommendations) == 0 {
        return nil, errors.New("未找到推荐商品")
    }

    // 计算订单总金额并创建订单详情
    var totalAmount float32 = 0
    var orderDetails []*userOrder.OrderDetail
    for _, rec := range recommendations {
        if err := validateRecommendation(&rec); err != nil {
            klog.Warnf("Invalid recommendation for userID %d: %+v, error: %v", req.UserId, rec, err)
            continue
        }
        amount := float32(rec.Price * float64(rec.Quantity))
        totalAmount += amount
        
        orderDetails = append(orderDetails, &userOrder.OrderDetail{
            ProductId: uint32(rec.ProductID),
            Name:      rec.Name,
            Number:    int32(rec.Quantity),
            Amount:    amount,
        })
    }

    if len(orderDetails) == 0 {
        return nil, errors.New("没有有效的商品可下单")
    }

    // 创建订单
    ctx, cancel := context.WithTimeout(ctx, s.timeout)
    defer cancel()
    
    // 使用Submit方法创建订单
    orderReq := &userOrder.OrderSubmitReq{
        UserId:    req.UserId,
        Amount:    totalAmount,
        PayMethod: 1, // 默认支付方式
        Remark:    "AI自动下单: " + req.Request,
        Order: &userOrder.OrderSubmitDetail{
            List: orderDetails,
        },
    }

    orderResp, err := s.orderClient.Submit(ctx, orderReq)
    if err != nil {
        klog.Errorf("Submit order failed for userID %d: %v", req.UserId, err)
        return nil, fmt.Errorf("创建订单失败: %w", err)
    }
    if orderResp == nil {
        return nil, errors.New("创建订单失败：服务返回为空")
    }

    // 返回订单ID
    orderID := strconv.FormatUint(uint64(orderResp.OrderId), 10)
    klog.Infof("Successfully created order %s for userID: %d", orderID, req.UserId)
    return &ai.AutoPlaceOrderResp{OrderId: orderID}, nil
}

func validateRecommendation(rec *aiutil.ProductRecommendation) error {
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