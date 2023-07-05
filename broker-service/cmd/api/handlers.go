package main

import (
	"broker/event"  // 导入自定义的 event 包
	"bytes"         // 导入 bytes 包用于操作字节
	"encoding/json" // 导入 json 包用于 JSON 编码和解码
	"errors"        // 导入 errors 包用于错误处理
	"net/http"      // 导入 net/http 包用于处理 HTTP 请求和响应
)

type RequestPayload struct { // 定义 RequestPayload 结构体，用于请求的载荷
	Action string      `json:"action"`         // 动作字段
	Auth   AuthPayload `json:"auth,omitempty"` // 认证字段，如果为空则省略
	Log    LogPayload  `json:"log,omitempty"`  // 日志字段，如果为空则省略
	Mail   MailPayload `json:"mail,omitempty"` // 邮件字段，如果为空则省略
}

type MailPayload struct { // 定义 MailPayload 结构体，用于邮件的载荷
	From    string `json:"from"`    // 发件人
	To      string `json:"to"`      // 收件人
	Subject string `json:"subject"` // 主题
	Message string `json:"message"` // 内容
}

type AuthPayload struct { // 定义 AuthPayload 结构体，用于认证的载荷
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

type LogPayload struct { // 定义 LogPayload 结构体，用于日志的载荷
	Name string `json:"name"` // 日志名称
	Data string `json:"data"` // 日志数据
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	// 声明一个名为 payload 的 jsonResponse 变量
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload) // 调用 writeJson 方法将 payload 写入响应
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload) // 从请求中读取 JSON 数据并解码到 requestPayload
	if err != nil {
		app.errorJson(w, err) // 如果解码出错，调用 errorJson 方法将错误信息写入响应
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth) // 调用 authenticate 方法进行认证
	case "log":
		app.logEventViaRabbit(w, requestPayload.Log) // 调用 logEventViaRabbit 方法发送日志事件
	case "mail":
		app.sendMail(w, requestPayload.Mail) // 调用 sendMail 方法发送邮件
	default:
		app.errorJson(w, errors.New("unknown action")) // 如果动作不匹配，则返回未知动作的错误
	}
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t") // 将邮件载荷 msg 转换为 JSON 格式的字节数据

	// 调用邮件服务
	mailServiceUrl := "http://mailer-service/send"

	// 发送 POST 请求到邮件服务
	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request) // 发送 HTTP 请求
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	// 确保返回的状态码正确
	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling mail service")) // 如果状态码不是 202 Accepted，则返回调用邮件服务时的错误
		return
	}

	// 发送 JSON 响应
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message send to " + msg.To

	app.writeJson(w, http.StatusAccepted, payload) // 调用 writeJson 方法将 payload 写入响应
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t") // 将日志载荷 entry 转换为 JSON 格式的字节数据

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData)) // 创建 POST 请求
	if err != nil {
		app.errorJson(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	response, err := client.Do(request) // 发送 HTTP 请求
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJson(w, http.StatusAccepted, payload) // 调用 writeJson 方法将 payload 写入响应
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t") // 将认证载荷 a 转换为 JSON 格式的字节数据

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData)) // 创建 POST 请求
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request) // 发送 HTTP 请求
	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()
	// 确保返回的状态码正确
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials")) // 如果状态码是 401 Unauthorized，则返回无效凭据的错误
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling auth service")) // 如果状态码不是 202 Accepted，则返回调用认证服务时的错误
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService) // 从响应中解码 JSON 数据到 jsonFromService
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload) // 调用 writeJson 方法将 payload 写入响应
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data) // 将日志事件推送到队列中
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJson(w, http.StatusAccepted, payload) // 调用 writeJson 方法将 payload 写入响应
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit) // 创建事件发射器
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t") // 将载荷转换为 JSON 格式的字节数据
	err = emitter.Push(string(j), "log.INFO")      // 将数据推送到队列中
	if err != nil {
		return err
	}
	return nil
}
