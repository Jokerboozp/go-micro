package main

import (
	"log-service/data"
	"net/http"
)

/**
该代码定义了一个 JSONPayload 结构体，用于表示 JSON 载荷的结构。
WriteLog 方法是处理写入日志请求的函数。
该函数首先将请求中的 JSON 数据解析为 JSONPayload 变量。
然后，它创建一个 data.LogEntry 对象，将解析后的 JSON 数据赋值给 Name 和 Data 字段。
接下来，它通过调用 app.Models.LogEntry.Insert 方法将日志条目插入数据库中。
如果插入过程中发生错误，它会调用 app.errorJson 方法返回错误的 JSON 响应。
最后，它创建一个 jsonResponse 对象作为成功的响应，并调用 app.writeJson 方法将响应写回客户端，状态码为 202（Accepted）。
这段代码主要用于处理写入日志的请求，它使用了自定义的 data 包来处理数据的插入，并依赖于外部的 http 包来处理 HTTP 请求和响应。
*/

// JSONPayload 结构体定义了 JSON 的载荷结构
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog 方法处理写入日志的请求
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// 将 JSON 解析为 JSONPayload 变量
	var requestPayload JSONPayload
	_ = app.readJson(w, r, &requestPayload)

	// 插入数据
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJson(w, http.StatusAccepted, resp)
}
