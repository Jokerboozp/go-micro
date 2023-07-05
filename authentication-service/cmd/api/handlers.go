package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

/**
该代码主要包含了两个函数：

Authenticate 函数用于处理认证请求。首先从请求中解析出邮箱和密码，然后根据邮箱从数据库中获取用户信息。接下来，验证密码是否匹配。如果认证成功，记录认证日志，并返回成功认证的响应。如果认证失败，返回相应的错误信息。
logRequest 函数用于发送日志请求给日志服务。首先构建日志数据结构，然后将数据转换为 JSON 格式。接着，发送 POST 请求给日志服务，并将日志数据作为请求体发送。如果发送过程中出现错误，则返回错误信息。
总结：该代码是一个认证服务，用于处理认证请求。其中，Authenticate 函数处理认证请求，验证用户的邮箱和密码，并记录认证日志。logRequest 函数负责发送日志请求给日志服务。
*/

// 处理认证请求
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// 定义请求的数据结构
	var requestPayload struct {
		Email    string `json:"email"`    // 邮箱
		Password string `json:"password"` // 密码
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest) // 返回错误的请求
		return
	}

	// 根据邮箱从数据库中获取用户
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest) // 返回无效凭证的错误
		return
	}

	// 验证密码是否匹配
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("invalid credentials!"), http.StatusUnauthorized) // 返回无效凭证的错误
		return
	}

	// 记录认证日志
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJson(w, err) // 返回日志记录错误
		return
	}

	// 构建响应数据
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload) // 返回成功认证的响应

}

// 记录日志请求
func (app *Config) logRequest(name, data string) error {
	// 定义日志数据结构
	var entry struct {
		Name string `json:"name"` // 日志名称
		Data string `json:"data"` // 日志数据
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://logger-service/log"

	// 发送日志请求
	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
