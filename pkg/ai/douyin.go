package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/123508/douyinshop/pkg/config"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type DouyinAI struct {
	client *arkruntime.Client
}

// 商品推荐请求的结构
type ProductRecommendation struct {
	ProductID   int64   `json:"product_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Confidence  float64 `json:"confidence"`
}

// 初始化豆包模型客户端
func NewDouyinAI() *DouyinAI {
// pkg/ai/douyin.go 中应该使用配置的超时时间
client := arkruntime.NewClientWithApiKey(
    config.Conf.VolcengineConfig.ApiKey,
    arkruntime.WithTimeout(time.Duration(config.Conf.VolcengineConfig.Timeout) * time.Second),
)
	
	return &DouyinAI{
		client: client,
	}
}

// Close 关闭资源
func (d *DouyinAI) Close() error {
	// 如果需要关闭资源，在这里实现
	return nil
}

// 分析用户订单需求
func (d *DouyinAI) AnalyzeOrderRequest(userRequest string) ([]ProductRecommendation, error) {
	// 构建提示词
	prompt := fmt.Sprintf(`作为一个电商智能助手，请分析以下用户需求并推荐合适的商品：
用户需求：%s

请按照以下JSON格式返回推荐商品：
{
    "products": [
        {
            "product_id": 商品ID,
            "name": "商品名称",
            "price": 价格,
            "quantity": 建议购买数量,
            "confidence": 推荐置信度
        }
    ]
}`, userRequest)

	ctx := context.Background()
	// 构建聊天请求
	req := model.ChatCompletionRequest{
		Model: config.Conf.VolcengineConfig.DouyinModel,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(prompt),
				},
			},
		},
		// 设置温度为0.7，增加一些创造性
		Temperature: volcengine.Float64(0.7),
	}

	// 调用模型
	resp, err := d.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("调用豆包模型失败: %v", err)
	}

	// 获取模型响应
	aiResponse := *resp.Choices[0].Message.Content.StringValue

	// 如果有推理内容，记录日志
	if resp.Choices[0].Message.ReasoningContent != nil {
		// TODO: 可以记录推理过程到日志系统
		_ = *resp.Choices[0].Message.ReasoningContent
	}

	// 解析模型返回的JSON
	var result struct {
		Products []ProductRecommendation `json:"products"`
	}
	if err := json.Unmarshal([]byte(aiResponse), &result); err != nil {
		return nil, fmt.Errorf("解析模型返回结果失败: %v", err)
	}

	return result.Products, nil
}

// 查询订单详情的格式化
func (d *DouyinAI) FormatOrderDetails(orderDetails map[string]interface{}) (string, error) {
	// 将订单信息转换为字符串
	orderJSON, err := json.Marshal(orderDetails)
	if err != nil {
		return "", fmt.Errorf("订单信息序列化失败: %v", err)
	}

	prompt := fmt.Sprintf(`请将以下订单信息格式化为易读的中文描述：
%s

请包含以下信息：
1. 订单基本信息（订单号、创建时间等）
2. 商品信息（名称、数量、价格等）
3. 订单状态
4. 支付信息
`, string(orderJSON))

	ctx := context.Background()
	req := model.ChatCompletionRequest{
		Model: config.Conf.VolcengineConfig.DouyinModel,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(prompt),
				},
			},
		},
		// 设置较低的温度，保持输出的一致性
		Temperature: volcengine.Float64(0.3),
	}

	resp, err := d.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("调用豆包模型失败: %v", err)
	}

	// 如果有推理内容，记录日志
	if resp.Choices[0].Message.ReasoningContent != nil {
		// TODO: 可以记录推理过程到日志系统
		_ = *resp.Choices[0].Message.ReasoningContent
	}

	return *resp.Choices[0].Message.Content.StringValue, nil
}

// 添加配置验证函数
func ValidateConfig(conf *config.VolcengineConfig) error {
	if conf.ApiKey == "" {
		return fmt.Errorf("缺少 API Key 配置")
	}
	if conf.DouyinModel == "" {
		return fmt.Errorf("缺少豆包模型配置")
	}
	return nil
}