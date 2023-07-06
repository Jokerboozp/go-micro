package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
该代码定义了一个 Config 结构体，用于存储应用程序的配置信息和数据模型。
全局变量 client 用于保存 MongoDB 的客户端连接。
main 函数是程序的入口函数。
在 main 函数中，首先连接到 MongoDB 数据库。
然后创建一个上下文以用于断开连接，并在函数结束时断开数据库连接。
创建 Config 对象，并将 MongoDB 客户端传递给 data.New 方法来创建数据模型。
启动 Web 服务器，监听指定的端口，并使用 app.routes() 方法作为处理程序。
如果启动服务器时发生错误，则打印错误信息。
connectToMongo 函数用于连接到 MongoDB 数据库。
在该函数中，创建连接选项，并设置认证信息。
然后使用连接选项连接到数据库。
如果连接过程中发生错误，则打印错误信息并返回错误。
如果连接成功，则打印连接成功的消息，并返回 MongoDB 客户端连接对象。
在程序中，使用 MongoDB 客户端连接进行数据库操作和数据访问。
*/

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// register the rpc server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	// start web server
	fmt.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (app *Config) rpcListen() error {
	fmt.Println("Starting rpc server on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			return err
		}

		go rpc.ServeConn(rpcConn)
	}
}

// connectToMongo 方法用于连接到 MongoDB 数据库
func connectToMongo() (*mongo.Client, error) {
	// 创建连接选项
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// 连接到数据库
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
