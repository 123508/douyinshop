package main

import (
	"context"
	"encoding/json"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

// OrderUserServiceImpl implements the last service interface defined in the IDL.
type OrderUserServiceImpl struct{}

const (
	serviceName = "order"
	goods       = "goods"
)

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		util.LogError("打开MySQL连接失败", "connectWithMySQL", "", err)
	}
	return DB
}

var rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

// GetOrderInfo 查询订单信息
func GetOrderInfo(ctx context.Context, orderId uint32) (*models.Order, error) {
	// 查询订单信息（Order）基本信息
	var order models.Order

	// 如果没有找到对应订单，返回错误信息
	if err := DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderError
	}

	return &order, nil
}

func GetOrderInfoWithCache(ctx context.Context, orderId uint32) (*models.Order, error) {
	// 查询订单信息（Order）基本信息
	var order models.Order

	orderKey := util.TakeKey(goods, "order", orderId)

	//查询单条order缓存
	jsonData, _ := rds.Get(ctx, orderKey).Result()

	//查询缓存失败
	if err := json.Unmarshal([]byte(jsonData), &order); err != nil {
		// 如果没有找到对应订单，返回错误信息
		if err := DB.Where("id = ?", orderId).First(&order).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderError
		}

		jsonData, err := json.Marshal(&order)

		if err != nil {
			util.LogError("序列化商品数据错误", "GetOrderInfoWithCache", "", err)
		} else {
			//将数据加入缓存
			if setErr := rds.Set(ctx, orderKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "GetOrderInfoWithCache", "将order加入缓存", setErr)
			}
		}

	} else {

		log.WithFields(log.Fields{
			"方法名": "GetOrderInfoWithCache",
		}).Info("查询order缓存成功")
	}

	return &order, nil
}

// GetOrderDetailsWithCache 查询指定订单的详细信息
func GetOrderDetailsWithCache(ctx context.Context, orderId uint32) ([]models.OrderDetail, error) {
	// 查询订单详情（List）
	orderDetails := make([]models.OrderDetail, 0)

	detailKey := util.TakeKey(goods, "orderDetail", orderId)

	//查询orderDetails缓存
	result, err := rds.Get(ctx, detailKey).Result()

	//查询orderDetail缓存失败
	if err != nil || json.Unmarshal([]byte(result), &orderDetails) != nil {

		// 如果查询订单详情失败，返回错误
		if dbErr := DB.Where("order_id = ?", orderId).Find(&orderDetails).Error; dbErr != nil {
			util.LogError("查询订单详情失败", "GetOrderDetailsWithCache", "查询订单详情", dbErr)
			return nil, SearchOrderDetailsError
		}

		jsonData, err := json.Marshal(&orderDetails)

		if err != nil {
			util.LogError("序列化数据失败", "GetOrderDetailsWithCache", "序列化orderDetail", err)
		} else {
			if setErr := rds.Set(ctx, detailKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "GetOrderDetailsWithCache", "将orderDetail加入缓存", setErr)
			}
		}

	} else {
		log.WithFields(log.Fields{
			"方法名": "GetOrderDetailsWithCache",
		}).Info("查询orderDetail缓存成功")
	}

	return orderDetails, nil
}

// GetOrderLogsWithCache 查询指定订单的日志信息
func GetOrderLogsWithCache(ctx context.Context, orderId uint32) ([]models.OrderStatusLog, error) {
	//查询订单日志详情
	orderLogs := make([]models.OrderStatusLog, 0)

	logKey := util.TakeKey(goods, "orderLog", orderId)

	res, err := rds.Get(ctx, logKey).Result()

	//查询缓存失败
	if err != nil || json.Unmarshal([]byte(res), &orderLogs) != nil {
		//查询失败报错
		if err = DB.Where("order_id = ?", orderId).Find(&orderLogs).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderLogsError
		}

		jsonData, err := json.Marshal(&orderLogs)

		if err != nil {
			util.LogError("序列化数据失败", "GetOrderLogsWithCache", "序列化orderLog", err)
		} else {
			if setErr := rds.Set(ctx, logKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "GetOrderLogsWithCache", "将orderLog加入缓存", setErr)
			}

		}
	} else {
		log.WithFields(log.Fields{
			"方法名": "GetOrderLogsWithCache",
		}).Info("查询orderLog缓存成功")
	}

	return orderLogs, nil
}

// 清除缓存
func cleanCache(ctx context.Context, orderId uint32) {
	//清除缓存
	util.CleanCache(rds, ctx, util.TakeKey(goods, "order", orderId))
	util.CleanCache(rds, ctx, util.TakeKey(goods, "orderDetail", orderId))
	util.CleanCache(rds, ctx, util.TakeKey(goods, "orderLog", orderId))
}

var SubmitOrderError = &errorno.BasicMessageError{Code: 500, Message: "订单提交失败"}

var SearchOrderError = &errorno.BasicMessageError{Code: 404, Message: "查询订单错误"}

var SearchOrderDetailsError = &errorno.BasicMessageError{Code: 500, Message: "查询订单详细信息错误"}

var SearchOrderLogsError = &errorno.BasicMessageError{Code: 404, Message: "查询订单日志错误"}

var CancelOrderError = &errorno.BasicMessageError{Code: 500, Message: "取消订单失败"}

var UnableChangeStatusError = &errorno.BasicMessageError{Code: 400, Message: "更新状态失败,该状态不允许被更新"}

var NotPayReminderError = &errorno.BasicMessageError{Code: 400, Message: "没有支付订单,无法提醒发货"}

var DeliveredReminderError = &errorno.BasicMessageError{Code: 400, Message: "商家已经发货,无需提醒"}

var CancelReminderError = &errorno.BasicMessageError{Code: 400, Message: "订单已经取消,无法提醒发货"}

var RefundingReminderError = &errorno.BasicMessageError{Code: 400, Message: "退款中,无法提醒发货"}

var RefundedReminderError = &errorno.BasicMessageError{Code: 400, Message: "已经退款,无法提醒发货"}

var RejectionReminderError = &errorno.BasicMessageError{Code: 400, Message: "商家拒绝发货,无法提醒"}

var StatusError = &errorno.BasicMessageError{Code: 400, Message: "不允许的行为"}

var CompletionOrderError = &errorno.BasicMessageError{Code: 400, Message: "无法确认收货"}

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}

// Submit implements the OrderUserServiceImpl interface.
// 用户提交订单
func (s *OrderUserServiceImpl) Submit(ctx context.Context, req *userOrder.OrderSubmitReq) (resp *userOrder.OrderSubmitResp, err error) {

	// 创建订单对象
	var order models.Order

	err = DB.Transaction(func(tx *gorm.DB) error {

		order.Number = uuid.New().String() // 使用 UUID 生成唯一的订单号
		order.UserId = ctx.Value("userId").(uint32)
		order.AddressBookId = req.AddressBookId
		order.PayMethod = req.PayMethod
		order.Remark = req.Remark
		order.ShopId = req.Order.ShopId
		order.FinalStatus = 0
		order.FinalVersion = 0

		var total float32

		for _, detail := range req.Order.List {
			total += detail.Amount * float32(detail.Number)
		}

		order.Amount = total

		if err = DB.Create(&order).Error; err != nil {
			return err
		}

		var orderDetails []models.OrderDetail
		for _, detail := range req.Order.List {

			orderDetail := models.OrderDetail{
				Name:      detail.Name,
				Image:     detail.Image,
				OrderId:   order.ID,
				ProductId: detail.ProductId,
				Number:    detail.Number,
				Amount:    detail.Amount,
			}

			// 将订单详情添加到列表中
			orderDetails = append(orderDetails, orderDetail)
		}

		//如果创建失败就返回错误
		if err = DB.Create(&orderDetails).Error; err != nil {
			return err
		}

		current := time.Now()

		task := time.Now().Add(15 * time.Minute)

		orderStatusLog := models.OrderStatusLog{
			OrderId:     order.ID,
			Status:      0, // 初始状态为待付款
			StartTime:   &current,
			EndTime:     &task,
			Description: "订单创建，待付款",
			Version:     0,
		}

		// 保存状态日志
		if err = DB.Create(&orderStatusLog).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return nil, SubmitOrderError
	}

	// 构建返回对象
	resp = &userOrder.OrderSubmitResp{
		OrderId:     order.ID,
		Number:      order.Number,
		OrderAmount: order.Amount,
	}

	return resp, nil
}

// History implements the OrderUserServiceImpl interface.
// 查询用户的历史订单
func (s *OrderUserServiceImpl) History(ctx context.Context, req *userOrder.HistoryReq) (resp *userOrder.HistoryResp, err error) {

	// 检查分页参数
	if req.Page < 1 || req.PageSize < 1 {
		return nil, BadPageOrPageSize
	}

	// 查询用户订单数据
	var orders []models.Order

	//使用Md5作为条件,支持复杂扩容
	key := util.Md5Hash(util.TakeKey(serviceName, "list", req.UserId, req.Page, req.PageSize, req.Status))

	result, _ := rds.Get(ctx, key).Result()

	list := make([]uint32, 0)
	fail := make([]uint32, 0)
	orderMap := make(map[uint32]models.Order)

	//查询id数组
	ok := json.Unmarshal([]byte(result), &list)

	//按照list查询详情缓存
	for _, v := range list {

		t := models.Order{}
		key := util.TakeKey(goods, "order", v)
		jsonData, err := rds.Get(ctx, key).Result()
		if err != nil || json.Unmarshal([]byte(jsonData), &t) != nil {
			fail = append(fail, v)
			continue
		}
		orderMap[v] = t
	}

	//计算缓存失效比率
	rate := 100
	if len(list) != 0 {
		rate = len(fail) * 100 / len(list)
	}

	if ok != nil || len(list) == 0 || rate > 30 {

		offset := (req.Page - 1) * req.PageSize

		err = DB.Where("user_id = ?", req.UserId).Offset(int(offset)).
			Limit(int(req.PageSize)).
			Find(&orders).Error

		if err != nil {
			log.Println(err)
			util.LogError("查询数据库错误", "History", "", err)
			return nil, SearchOrderError
		}

		//构建订单ID数组
		idList := make([]uint32, 0, len(orders))
		for _, v := range orders {
			idList = append(idList, v.ID)
		}

		//序列化并存储ID列表缓存
		jsonData, err := json.Marshal(&idList)

		if err != nil {
			util.LogError("序列化订单数组错误", "History", "请求hash为"+key, err)
		} else {
			// 设置随机过期时间，3~6 分钟
			if setErr := rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("存入order缓存失败", "History", "", setErr)
			}
		}

		//订单详情分级存储
		for _, v := range orders {
			key := util.TakeKey(goods, "order", v.ID)
			jsonData, err := json.Marshal(&v)
			if err != nil {
				util.LogError("序列化订单错误", "History", "order的id为"+strconv.Itoa(int(v.ID)), err)
			} else {
				if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
					util.LogError("缓存订单失败", "History", "order的id为"+strconv.Itoa(int(v.ID)), err)
				}
			}
		}

	} else {
		//缓存失效比例较低,逐条查询并放入缓存
		if len(fail) > 1 {
			var missedOrders []models.Order
			if err := DB.Where("id IN ?", fail).Find(&missedOrders).Error; err != nil {
				util.LogError("查询order错误", "History", "", err)
				return nil, SearchOrderError
			}
			for _, o := range missedOrders {
				orderMap[o.ID] = o
				//放入缓存
				key := util.TakeKey(goods, "order", o.ID)

				jsonData, err := json.Marshal(&o)
				if err != nil {
					util.LogError("序列化订单错误", "History", "订单id为"+strconv.Itoa(int(o.ID)), err)
				} else {
					if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
						util.LogError("缓存订单失败", "History", "订单id为"+strconv.Itoa(int(o.ID)), err)
					}
				}
			}
		} else if len(fail) == 1 {
			t := models.Order{}

			//查询异常,返回错误
			if err := DB.Where("id = ?", fail[0]).First(&t).Error; err != nil {
				util.LogError("查询order错误", "History", "", err)
				return nil, SearchOrderError
			}

			orderMap[fail[0]] = t

			//放入缓存
			key := util.TakeKey(goods, "order", fail[0])

			jsonData, err := json.Marshal(&t)
			if err != nil {
				util.LogError("序列化订单错误", "History", "订单id为"+strconv.Itoa(int(fail[0])), err)
			} else {
				if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
					util.LogError("缓存订单失败", "History", "订单id为"+strconv.Itoa(int(fail[0])), err)
				}
			}
		}

		// 最终按list顺序组装orders
		orders = make([]models.Order, 0, len(list))
		for _, id := range list {
			if v, ok := orderMap[id]; ok {
				orders = append(orders, v)
			}
		}

		log.WithFields(log.Fields{
			"方法名": "History",
		}).Info("查询order缓存成功")
	}

	// 构建订单响应数据
	orderList := make([]*order_common.OrderResp, len(orders))
	for i, order := range orders {
		// 构建订单响应
		orderResp := &order_common.OrderResp{
			Order: &order_common.Order{
				ID:          order.ID,
				Number:      order.Number,
				Amount:      order.Amount,
				FinalStatus: order.FinalStatus,
			},
		}
		orderList[i] = orderResp
	}

	// 返回响应
	return &userOrder.HistoryResp{
		List:     orderList,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    uint32(len(orders)),
	}, nil
}

// Detail implements the OrderUserServiceImpl interface.
// 查询订单的详细信息
func (s *OrderUserServiceImpl) Detail(ctx context.Context, req *order_common.OrderReq) (resp *order_common.OrderResp, err error) {

	order, err := GetOrderInfoWithCache(ctx, req.OrderId)

	if err != nil {
		return nil, err
	}

	orderDetails, err := GetOrderDetailsWithCache(ctx, req.OrderId)

	if err != nil {
		return nil, err
	}

	orderLogs, err := GetOrderLogsWithCache(ctx, req.OrderId)

	if err != nil {
		return nil, err
	}

	var statusList []*order_common.Status

	for _, k := range orderLogs {
		statusList = append(statusList, &order_common.Status{
			StartTime:   k.StartTime.String(),
			Status:      k.Status,
			EndTime:     k.EndTime.String(),
			Description: k.Description,
		})
	}

	// 处理类型
	orderCommon := &order_common.Order{
		Number:        order.Number,
		UserId:        order.UserId,
		PayMethod:     order.PayMethod,
		AddressBookId: order.AddressBookId,
		Amount:        order.Amount,
		Remark:        order.Remark,
		Phone:         order.Phone,
		Address:       order.Address,
		Username:      order.UserName,
		Consignee:     order.Consignee,
		Status:        statusList,
	}

	var orderDetail []*order_common.OrderDetail

	for _, k := range orderDetails {
		orderDetail = append(orderDetail, &order_common.OrderDetail{
			Name:      k.Name,
			Image:     k.Image,
			OrderId:   k.OrderId,
			ProductId: k.ProductId,
			Number:    k.Number,
			Amount:    k.Amount,
		})
	}

	// 构建并返回响应数据
	resp = &order_common.OrderResp{
		Order:        orderCommon,
		OrderDetails: orderDetail,
	}

	return resp, nil
}

// Cancel implements the OrderUserServiceImpl interface.
// 取消订单
// 步骤:查询要取消的订单->查询不到，返回错误，否则继续->判断该订单的status是否<5,否返回错误,是继续->修改订单日志存储状态，如果报错就返回错误->返回正确响应
func (s *OrderUserServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

	defer cleanCache(ctx, req.OrderId)

	//查询要取消的订单
	order, err := GetOrderInfo(ctx, req.OrderId)

	//查询要取消的订单失败
	if err != nil {
		return nil, err
	}

	//查询订单状态
	var status models.OrderStatusLog

	//查询订单日志异常
	if err = DB.Where("order_id = ? and version = ?", req.OrderId, order.FinalVersion).Last(&status).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogsError
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	//订单状态错误
	if order.FinalStatus > 5 {
		return nil, UnableChangeStatusError
	}

	//决定下一个状态
	var Status uint32

	var Description string

	if order.FinalStatus == 0 {
		Status = 6
		Description = "已取消"
	} else {
		Status = 7
		Description = "退款中"
	}

	//创建新状态
	newStatus := models.OrderStatusLog{
		StartTime:   &currentTime,
		Status:      Status,
		EndTime:     nil,
		Description: Description,
		Version:     status.Version + 1,
	}

	// 更新订单详情的状态为“已取消”并记录到 OrderStatusLog
	err = DB.Transaction(func(tx *gorm.DB) error {

		//将原有状态的结束时间修改
		if err = DB.Model(&models.OrderStatusLog{}).Where("id = ?", status.ID).Update("end_time", currentTime).Error; err != nil {
			return err
		}

		//插入新的状态
		if err = DB.Create(&newStatus).Error; err != nil {
			return err
		}

		//为修改订单为取消状态
		if err = DB.Model(&models.Order{}).Where("id = ?", req.OrderId).Update("final_status", newStatus.Status).Update("final_version", newStatus.Version).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return nil, CancelOrderError
	}

	return &order_common.Empty{}, nil
}

// Reminder implements the OrderUserServiceImpl interface.
// 提醒商家发货
// 步骤:查询订单是否属于用户->判断status==2,不是就返回错误,是继续(增加健壮性)->通过订单查询商家id,查不到就报错,否则继续->发送信息
func (s *OrderUserServiceImpl) Reminder(ctx context.Context, req *userOrder.ReminderReq) (resp *order_common.Empty, err error) {

	defer cleanCache(ctx, req.OrderId)

	// 查询该用户的订单信息
	var order models.Order

	//查询订单失败,返回异常
	if err = DB.Where("id = ? AND user_id = ?", req.OrderId, ctx.Value("userId").(uint32)).First(&order).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderError
	}

	//如果状态不是商家已接单就返回错误
	switch order.FinalStatus {
	case 0:
		return nil, NotPayReminderError
	case 1, 2: //只允许这两个状态通行
	case 3, 4, 5:
		return nil, DeliveredReminderError
	case 6:
		return nil, CancelReminderError
	case 7:
		return nil, RefundingReminderError
	case 8:
		return nil, RefundedReminderError
	case 9:
		return nil, RejectionReminderError
	default:
		return nil, StatusError
	}

	//提醒商家发货逻辑
	err = util.SendMessage("order.direct", "message", "orderId:"+strconv.Itoa(int(req.OrderId))+",shopId:"+strconv.Itoa(int(order.ShopId)), 1)

	if err != nil {
		return nil, err
	}

	return &order_common.Empty{}, nil
}

// Complete implements the OrderUserServiceImpl interface.
// 确认收货
// 订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单 取消退款(直接回到上一步即可)
// 步骤:查询要完成的订单->查询不到，返回错误，否则继续->判断该订单的status是否为4,否返回错误,是继续->修改订单日志存储状态，如果报错就返回错误->返回正确响应
func (s *OrderUserServiceImpl) Complete(ctx context.Context, req *userOrder.CompleteReq) (resp *order_common.Empty, err error) {

	defer cleanCache(ctx, req.OrderId)

	// 查询订单信息
	order, err := GetOrderInfo(ctx, req.OrderId)

	//查询不到订单信息,返回异常
	if err != nil {
		return nil, err
	}

	//查询订单状态信息
	var status models.OrderStatusLog

	//查询订单日志失败,返回异常
	if err = DB.Where("order_id = ? and version = ?", req.OrderId, order.FinalVersion).Last(&status).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogsError
	}

	//如果订单状态不为待派送就无法收货
	if order.FinalStatus != 4 && order.FinalStatus != 3 {
		return nil, UnableChangeStatusError
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	//创建完成状态
	newStatus := models.OrderStatusLog{
		StartTime: &currentTime,
		Status:    5,
		EndTime:   nil,
		Version:   status.Version + 1,
	}
	err = DB.Transaction(func(tx *gorm.DB) error {

		//将原有状态的结束时间修改
		if err = DB.Model(&models.OrderStatusLog{}).Where("id = ?", status.ID).Update("end_time", currentTime).Error; err != nil {
			return err
		}

		//插入新的状态
		if err = DB.Create(&newStatus).Error; err != nil {
			return err
		}

		//为修改订单为完成状态
		if err = DB.Model(&models.Order{}).Where("id = ?", req.OrderId).Update("final_status", newStatus.Status).Update("final_version", newStatus.Version).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return nil, CompletionOrderError
	}

	return &order_common.Empty{}, nil
}
