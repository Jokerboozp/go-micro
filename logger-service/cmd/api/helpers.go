package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

/**
该代码定义了一个 jsonResponse 结构体，用于表示 JSON 响应的结构。
readJson 方法用于解析请求中的 JSON 数据。
在 readJson 方法中，它设置了一个最大字节数限制，然后使用 json.Decoder 将请求体解码为指定的数据结构。
该方法还检查请求体是否只包含一个 JSON 值。
writeJson 方法用于将数据以 JSON 格式写入响应。
它将提供的数据转换为 JSON 格式，并将其写入 ResponseWriter 中。
如果提供了 headers 参数，则将其中的键值对设置到 ResponseWriter 的 Header 中。
errorJson 方法用于返回错误的 JSON 响应。
它根据提供的错误信息和状态码创建一个包含错误信息的 jsonResponse 对象，并调用 writeJson 方法将其作为响应写入 ResponseWriter。
这些方法主要用于处理 JSON 数据的读取、写入和错误响应的生成，以便在 HTTP 请求和响应过程中进行数据的序列化和反序列化。
*/

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJson 方法用于解析请求中的 JSON 数据
func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// 检查请求体是否只有一个 JSON 值
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

// writeJson 方法用于将数据以 JSON 格式写入响应
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

// errorJson 方法用于返回错误的 JSON 响应
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
