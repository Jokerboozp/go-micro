package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// reading json
func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //one megabyte
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

// write json
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
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

// error json
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
