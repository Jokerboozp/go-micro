package main

import (
	"fmt"
	"net/http"
)

/**
该代码定义了一个 SendMail 方法，用于处理发送邮件的请求。
方法接收一个 http.ResponseWriter 和一个 http.Request 作为参数。
该方法首先定义了一个 mailMessage 结构体，表示邮件消息的字段。
然后，创建一个 requestPayload 变量，用于接收请求中的 JSON 数据。
通过调用 app.readJson 方法将请求中的 JSON 数据解码并存储到 requestPayload 中。
如果解码过程中发生错误，则打印错误信息并返回错误的 JSON 响应。
创建一个 Message 对象，使用 requestPayload 的字段填充邮件消息。
调用 app.Mailer.SendSMTPMessage 方法发送邮件消息。
如果发送过程中发生错误，则打印错误信息并返回错误的 JSON 响应。
创建一个 payload 对象，表示成功发送邮件的 JSON 响应。
调用 app.writeJson 方法将 JSON 响应写入 http.ResponseWriter 中。
*/

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	// 定义邮件消息结构体
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	// 创建 mailMessage 对象
	var requestPayload mailMessage

	// 读取 JSON 数据到 requestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		fmt.Println(err)
		app.errorJson(w, err)
		return
	}

	// 创建邮件消息对象
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// 发送邮件消息
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		fmt.Println(err)
		app.errorJson(w, err)
		return
	}

	// 创建响应数据对象
	payload := jsonResponse{
		Error:   false,
		Message: "Sent to " + requestPayload.To,
	}

	// 写入 JSON 响应
	app.writeJson(w, http.StatusAccepted, payload)
}
