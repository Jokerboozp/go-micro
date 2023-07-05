package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

/**
该代码是一个身份验证服务的主程序。主要包括以下功能：

连接到 PostgreSQL 数据库。
启动 HTTP 服务器并监听指定的端口。
处理数据库连接的重试逻辑。
总结：该程序启动一个身份验证服务，它通过连接到 PostgreSQL 数据库来提供身份验证功能。在连接数据库时，使用了重试逻辑，如果数据库尚未就绪，则会进行重试。一旦连接成功，程序将启动 HTTP 服务器，并使用指定的端口提供身份验证服务。
*/

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// 连接到数据库
	conn := connectToDB()
	if conn == nil {
		log.Panic("无法连接到 Postgres 数据库！")
	}

	// 设置配置
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// 打开数据库连接
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 连接到数据库
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres 数据库尚未就绪...")
			counts++
		} else {
			log.Println("已连接到 Postgres 数据库！")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("等待两秒钟...")
		time.Sleep(2 * time.Second)
		continue
	}
}
