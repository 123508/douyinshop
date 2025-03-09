package util

import (
	"fmt"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)

const (
	Direct  string = "direct"
	Fanout  string = "fanout"
	Topic   string = "topic"
	Headers string = "headers"
)

// DeclareAndBind 使用给定的绑定规范来声明队列、交换机并进行绑定
func DeclareAndBind(amqpChan *amqp.Channel, exchangeName, exchangeType, queueName, routingKey string) error {
	// 声明交换机
	err := amqpChan.ExchangeDeclare(
		exchangeName, // 交换机名称
		exchangeType, // 交换机类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部
		false,        // 不等待服务器确认
		nil,          // 额外参数
	)
	if err != nil {
		return fmt.Errorf("交换机声明失败,原因: %v", err)
	}

	// 声明队列
	queue, err := amqpChan.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 排他
		false,     // 不自动删除
		false,     // 不等待服务器确认
		nil,       // 额外参数
	)
	if err != nil {
		return fmt.Errorf("队列声明失败,原因: %v", err)
	}

	// 绑定队列到交换机
	err = amqpChan.QueueBind(
		queue.Name,   // 队列名称
		routingKey,   // 路由键
		exchangeName, // 交换机名称
		false,        // 不等待服务器确认
		nil,          // 额外参数
	)
	if err != nil {
		return fmt.Errorf("队列绑定失败,原因: %v", err)
	}

	fmt.Printf("声明并绑定队列'%s'到交换机'%s',绑定的路由键为:'%s'\n",
		queue.Name, exchangeName, routingKey)
	return nil
}

// SendMessage 发送的消息如果没有对应的队列接收,会直接丢弃消息,没有添加生产者确认机制
func SendMessage(exchange, routingKey, message string, number int) error {
	// 连接到RabbitMQ服务器
	conn, err := amqp.Dial("amqp://" +
		config.Conf.RabbitmqConfig.Username +
		":" +
		config.Conf.RabbitmqConfig.Password +
		"@" +
		config.Conf.RabbitmqConfig.Host +
		":" +
		strconv.Itoa(config.Conf.RabbitmqConfig.Port) +
		config.Conf.RabbitmqConfig.VirtualHost)

	if err != nil {
		return fmt.Errorf("无法连接到RabbitMQ,原因: %v", err)
	}
	defer conn.Close()

	// 打开一个通道
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("无法打开通道,原因: %v", err)
	}
	defer ch.Close()

	//健壮性判断
	if number <= 0 {
		number = 1
	}

	//多次发送
	for i := 0; i < number; i++ {
		// 发送消息到队列
		err = ch.Publish(
			exchange,   // 交换机名称
			routingKey, // 路由键，这里使用队列名称
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			})
		if err != nil {
			return fmt.Errorf("发送消息失败,原因: %v", err)
		}
	}

	log.Printf("发送消息成功,共计%d条,消息内容为%s\n", number, message)
	return nil
}

// Handler 请给出一个处理函数
type Handler interface {
	HandlerMessage(msg amqp.Delivery)
}

func ReceiveMessages(exchangeName, exchangeType, queueName, routingKey string, t Handler) {

	stopChan := make(chan bool)

	go func() {
		for {
			select {
			case <-stopChan:
				fmt.Println("协程退出")
				return // 退出协程
			default:

				count := 0

				conn, err := amqp.Dial("amqp://" +
					config.Conf.RabbitmqConfig.Username +
					":" +
					config.Conf.RabbitmqConfig.Password +
					"@" +
					config.Conf.RabbitmqConfig.Host +
					":" +
					strconv.Itoa(config.Conf.RabbitmqConfig.Port) +
					config.Conf.RabbitmqConfig.VirtualHost)

				if err != nil {
					log.Printf("无法连接到rabbitmq: %v", err)
					time.Sleep(5 * time.Second) // 等待一段时间后重连

					//重试超过三次退出
					if count++; count > 3 {
						conn.Close()
						//向停止通道发送消息,退出协程
						stopChan <- true
					}

					continue
				} else {
					fmt.Println("连接到rabbitmq成功...")
				}

				count = 0

				ch, err := conn.Channel()
				if err != nil {
					log.Printf("无法打开连接通道: %v", err)
					conn.Close()
					time.Sleep(5 * time.Second) // 等待一段时间后重连

					//重试超过三次退出
					if count++; count > 3 {
						ch.Close()
						conn.Close()
						//向停止通道发送消息,退出协程
						stopChan <- true
					}
					continue
				} else {
					fmt.Println("打开连接通道成功...")
				}

				err = DeclareAndBind(ch, exchangeName, exchangeType, queueName, routingKey)

				if err != nil {
					log.Printf("无法绑定交换机和队列:%v", err)
					ch.Close()
					conn.Close()
					stopChan <- true
				}

				q, err := ch.QueueDeclarePassive(queueName, true, false, false, false, nil)
				count = 0
				if err != nil {
					log.Printf("无法声明队列: %v", err)
					ch.Close()
					conn.Close()
					time.Sleep(5 * time.Second) // 等待一段时间后重连

					//重试超过三次退出
					if count++; count > 3 {
						ch.Close()
						conn.Close()
						//向停止通道发送消息,退出协程
						stopChan <- true
					}
					continue
				} else {
					fmt.Println("声明队列成功,队列为:", q.Name)
				}

				// 设置Qos，确保每次只预取一条消息
				err = ch.Qos(1, 0, false)
				count = 0
				if err != nil {
					log.Fatalf("无法设置Qos,原因: %v", err)

					//重试超过三次退出
					if count++; count > 3 {
						ch.Close()
						conn.Close()
						//向停止通道发送消息,退出协程
						stopChan <- true
					}
					continue
				} else {
					log.Println("设置Qos成功...")
				}

				messages, err := ch.Consume(
					q.Name, // queue name
					"",     // consumer
					true,   // auto-ack
					false,  // exclusive
					false,  // no-local
					false,  // no-wait
					nil,    // args
				)

				count = 0
				if err != nil {
					log.Printf("无法注册消费者,原因: %v", err)
					ch.Close()
					conn.Close()
					time.Sleep(5 * time.Second) // 等待一段时间后重连

					//重试超过三次退出
					if count++; count > 3 {
						ch.Close()
						conn.Close()
						//向停止通道发送消息,退出协程
						stopChan <- true
					}
					continue
				} else {
					fmt.Println("注册消费者成功...")
				}

				log.Println("开始接收消息...")

				for msg := range messages {

					if t == nil {
						log.Println("接受到的消息是:", string(msg.Body))
					} else {
						t.HandlerMessage(msg)
					}

				}

				// 如果退出循环，则关闭通道和连接
				ch.Close()
				conn.Close()
				log.Println("消费者被关闭,重连中...")
				time.Sleep(5 * time.Second) // 等待一段时间后重连
			}
		}
	}()
}
