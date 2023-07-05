package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

/**
该代码定义了一个 routes 方法，用于设置路由和中间件的 HTTP 处理程序。以下是代码的注释：

mux := chi.NewRouter() 创建一个新的路由器实例。
mux.Use(cors.Handler(cors.Options{...})) 设置 CORS（跨域资源共享）中间件，指定允许连接的来源、方法、头部等。
mux.Use(middleware.Heartbeat("/ping")) 添加心跳路由中间件，将 "/ping" 映射到一个简单的处理函数。
mux.Post("/authenticate", app.Authenticate) 将 "/authenticate" 路由映射到 app.Authenticate 方法。
该方法返回一个 http.Handler 对象，可以用于处理 HTTP 请求。
*/

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// 指定允许连接的来源
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 使用中间件添加心跳路由
	mux.Use(middleware.Heartbeat("/ping"))

	// 将 /authenticate 路由映射到 app.Authenticate 方法
	mux.Post("/authenticate", app.Authenticate)

	return mux
}
