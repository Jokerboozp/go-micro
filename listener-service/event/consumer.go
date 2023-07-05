package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

/**
该代码定义了一个 Consumer 结构体和一些与消费消息相关的函数。
NewConsumer 函数用于创建一个新的 Consumer 实例，它接收一个 amqp.Connection 参数，并返回一个 Consumer 实例和一个错误。
setup 函数用于设置 Consumer，它创建一个 AMQP 通道并调用 declareExchange 函数声明交换机。
Payload 结构体定义了消息的载荷，具有 Name 和 Data 字段。
Listen 函数用于监听消息队列中的消息，它创建一个随机队列，并将队列与指定的主题绑定。然后，它通过通道消费消息，并将每个消息的载荷传递给 handlePayload 函数进行处理。
handlePayload 函数根据消息的名称进行处理。如果名称为 "log" 或 "event"，则调用 logEvent 函数记录事件。否则，不做任何操作。
logEvent 函数将事件记录到远程日志服务。它将载荷转换为 JSON 格式，并使用 HTTP POST 请求将日志发送到日志服务的 URL。如果请求成功并返回状态码为 202（Accepted），则表示记录成功。否则，返回错误。
该代码使用了 RabbitMQ 的 amqp091-go 包和标准库中的 JSON 和 HTTP 包来实现与 RabbitMQ 和远程日志服务的交互。
*/

// Consumer 结构体用于消费消息
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer 函数创建一个新的 Consumer 实例
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup 函数用于设置 Consumer
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Payload 结构体定义消息的载荷
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen 函数用于监听消息队列中的消息
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange,Queue] [logs_topic,%s]\n", q.Name)
	<-forever

	return nil
}

// handlePayload 函数处理消息的载荷
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	case "auth":
	default:
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// logEvent 函数记录事件
func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
