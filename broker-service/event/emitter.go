package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

/**
该代码定义了一个名为 Emitter 的结构体，表示事件的发布者。
setup 方法用于设置 Emitter，其中会创建一个通道，并调用 declareExchange 方法来声明一个交换机。
Push 方法用于将事件推送到消息队列。它创建一个通道并使用 channel.Publish 方法发布消息到交换机。
NewEventEmitter 函数用于创建一个新的 Emitter 实例，并调用 setup 方法进行设置。
该代码与消息队列相关，使用了 RabbitMQ 的 amqp091-go 包来实现与 RabbitMQ 的交互。
*/

type Emitter struct {
	connection *amqp.Connection
}

// 设置 Emitter
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	return declareExchange(channel) // 声明一个交换机
}

// 推送事件到消息队列
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel") // 打印推送消息的提示信息

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event), // 设置消息的内容
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// 创建一个新的 Event Emitter 实例
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup() // 调用 setup 方法设置 Emitter
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
