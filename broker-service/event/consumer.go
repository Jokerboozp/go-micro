package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

/**
该代码定义了一个名为 Consumer 的结构体，表示消息队列的消费者。
NewConsumer 函数用于创建一个新的 Consumer 实例，并调用 setup 方法进行设置。
setup 方法使用连接创建一个通道，并调用 declareExchange 方法来声明一个交换机。
Payload 结构体定义了消息负载的格式。
Listen 方法用于监听消息队列，并处理接收到的消息。
handlePayload 函数根据负载的名称进行不同的处理，其中 "log" 和 "event" 调用 logEvent 方法记录日志。
logEvent 函数将负载转换为 JSON 格式的字节数据，并发送一个 HTTP POST 请求到日志服务的 URL。
*/

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// 创建一个新的 Consumer 实例
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup() // 调用 setup 方法设置消费者
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// 设置消费者
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel) // 调用声明交换机的方法
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// 监听消息队列
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch) // 声明一个随机队列
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

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil) // 消费消息队列中的消息
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload) // 解码消息体中的 JSON 数据到 payload 变量中

			go handlePayload(payload) // 异步处理消息的负载
		}
	}()

	fmt.Printf("Waiting for message [Exchange,Queue] [logs_topic,%s]\n", q.Name)
	<-forever // 永远等待，保持监听状态

	return nil
}

// 处理消息负载
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// 记录日志
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

// 记录日志
func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t") // 将 entry 对象转换为 JSON 格式的字节数据

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData)) // 创建一个 HTTP POST 请求
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json") // 设置请求的 Content-Type 为 application/json

	client := http.Client{}

	response, err := client.Do(request) // 发送请求
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
