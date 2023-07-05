package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

/**
该代码定义了一个 Config 结构体，用于存储邮件配置信息。
常量 webPort 设置为字符串 "80"，表示 web 服务器的端口号。
main 函数是程序的入口函数。
在 main 函数中，创建了一个 Config 对象 app，并初始化其 Mailer 字段为通过 createMail 函数创建的邮件配置对象。
输出日志信息，表示开始在指定端口上启动邮件服务。
创建一个 HTTP 服务器实例 srv，设置其监听地址为 webPort，处理器为 app.routes()。
使用 ListenAndServe 方法启动服务器，并处理请求。
如果发生错误，打印错误信息。
createMail 函数用于创建邮件配置对象。
通过读取环境变量获取邮件相关的配置信息，并构建 Mail 对象。
返回构建好的 Mail 对象。
*/

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {
	// 创建配置对象并初始化邮件配置
	app := Config{
		Mailer: createMail(),
	}
	log.Println("Start mail service on port", webPort)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// 启动服务器监听并处理请求
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

// 创建邮件配置对象
func createMail() Mail {
	// 从环境变量中获取邮件相关配置
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_FROMNAME"),
		FromAddress: os.Getenv("MAIL_FROMADDRESS"),
	}

	return m
}
