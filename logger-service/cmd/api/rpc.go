package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

/**
该代码定义了一个 RPCServer 结构体，表示 RPC 服务器。
该结构体没有任何字段，仅用于组织和管理相关方法。
该代码还定义了一个 RPCPayload 结构体，表示 RPC 传输的有效载荷。
LogInfo 是 RPCServer 结构体的方法，用于处理日志信息。
该方法接收一个 RPCPayload 对象作为输入参数和一个指向字符串指针的响应参数。
在该方法中，首先获取名为 "logs" 的数据库中的 "logs" 集合。
然后，将传入的日志有效载荷插入到集合中作为一条新的日志记录。
如果在插入过程中发生错误，则打印错误信息并返回错误。
最后，设置响应字符串，表示成功处理了 RPC 请求，然后返回 nil 表示没有发生错误。
*/

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

// LogInfo 是 RPCServer 结构体的方法，用于处理日志信息
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	// 获取名为 "logs" 的数据库中的 "logs" 集合
	collection := client.Database("logs").Collection("logs")

	// 向集合中插入一条日志记录
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error writing to mongo:", err)
		return err
	}

	// 设置响应字符串，表示成功处理了 RPC 请求
	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
