package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

/**
该代码定义了一个结构体 jsonResponse，表示 JSON 响应的格式，包括是否出错、消息和数据。
readJson 函数用于从请求中读取 JSON 数据，并将其解码到指定的变量中。
writeJson 函数用于将数据编码为 JSON，并写入响应中。
errorJson 函数用于生成一个出错的 JSON 响应，并将其写入响应中，同时设置相应的状态码。
*/

type jsonResponse struct {
	Error   bool   `json:"error"`          // 表示响应是否出错
	Message string `json:"message"`        // 响应的消息
	Data    any    `json:"data,omitempty"` // 响应的数据，可选
}

// 读取 JSON
func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 一兆字节
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data) // 解码请求体的 JSON 数据到 data 变量中
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value") // 确保请求体只有一个 JSON 值
	}
	return nil
}

// 写入 JSON
func (app *Config) writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data) // 将 data 转换为 JSON 格式的字节数据
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value // 将提供的 headers 中的键值对设置到 ResponseWriter 的 Header 中
		}
	}

	w.Header().Set("Content-Type", "application/json") // 设置 ResponseWriter 的 Content-Type 为 application/json
	w.WriteHeader(status)                              // 设置 ResponseWriter 的状态码为指定的 status

	_, err = w.Write(out) // 将 JSON 数据写入 ResponseWriter
	if err != nil {
		return err
	}

	return nil
}

// 错误 JSON
func (app *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJson(w, statusCode, payload) // 调用 writeJson 方法将错误信息写入响应
}
