# Go-Micro Udemy Lesson

## 1、Building a simple front end and one Microservice

### 1.1、front-end

- 略过，主要进行Go后端，前端直接下载课程模板

### 1.2、broker-service

- 创建broker-service文件夹
- 在broker-service文件夹中执行`go mod init broker`,创建`gomod`
- 创建cmd/api/main.go
- 安装`chi`模块,一个第三方路由包
```shell
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
go get github.com/go-chi/cors
```
- 创建project文件夹
- 在broker-service创建`broker-service.dockerfile`
- 在project中创建`docker-compose.yml`
- 在project文件夹路径下执行`docker-compose up -d`才能创建对应的docker容器并运行