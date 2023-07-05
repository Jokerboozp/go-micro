package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

/**
该代码提供了一些用于处理 JSON 数据的函数：

readJson 函数用于从请求中读取 JSON 数据并解析到指定的数据结构中。
writeJson 函数用于将数据转换为 JSON 格式，并写入到响应中。
errorJson 函数用于返回一个包含错误信息的 JSON 响应。该函数会设置错误标志为 true，并将错误消息作为响应的消息字段。
总结：该代码提供了一些通用的函数，用于处理 JSON 数据的读取、写入和错误响应。这些函数可以在处理 HTTP 请求时方便地进行 JSON 数据的解析和生成响应。
*/

type jsonResponse struct {
	Error   bool   `json:"error"`          // 错误标志
	Message string `json:"message"`        // 错误消息
	Data    any    `json:"data,omitempty"` // 数据
}

// 读取 JSON 数据
func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

// 写入 JSON 数据
func (app *Config) writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// 将 data 转换为 JSON 格式
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 如果提供了 headers 参数，则将其中的键值对设置到 ResponseWriter 的 Header 中
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// 设置 ResponseWriter 的 Content-Type 为 application/json
	w.Header().Set("Content-Type", "application/json")

	// 设置 ResponseWriter 的状态码为指定的 status
	w.WriteHeader(status)

	// 将 JSON 数据写入 ResponseWriter
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	// 返回 nil，表示没有发生错误
	return nil
}

// 返回错误的 JSON 响应
func (app *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJson(w, statusCode, payload)
}
