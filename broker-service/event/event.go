package event

import amqp "github.com/rabbitmq/amqp091-go"

/**
该代码定义了两个函数，用于在消息队列中声明交换机和队列。
declareExchange 函数用于声明一个交换机，它接收一个 amqp.Channel 作为参数，并使用 ch.ExchangeDeclare 方法进行声明。该函数指定了交换机的名称、类型以及是否持久化等属性。
declareRandomQueue 函数用于声明一个随机队列，它接收一个 amqp.Channel 作为参数，并使用 ch.QueueDeclare 方法进行声明。该函数指定了队列的名称、是否持久化、是否在使用时删除等属性。
这些函数使用了 RabbitMQ 的 amqp091-go 包来实现与 RabbitMQ 的交互。
*/

// 声明一个交换机
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // 交换机名称
		"topic",      // 交换机类型
		true,         // 是否持久化
		false,        // 是否自动删除
		false,        // 是否为内部交换机
		false,        // 是否等待
		nil,          // 其他参数
	)
}

// 声明一个随机队列
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // 队列名称
		false, // 是否持久化
		false, // 是否在使用时删除
		true,  // 是否独占
		false, // 是否等待
		nil,   // 其他参数
	)
}
