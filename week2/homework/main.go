package main

import (
	"fmt"
	"net/http"
	"os"
)

// 有两种方式实现http server
// 1. 通过 golang 自带的 net/http server 实现， net/http 包已经实现了 http包的解析
// 2. 通过 golang 自带的 net/tcp 去实现， 这种方式需要自己解析http请求，自己拼装http响应
func main() {
	srv := &http.Server{
		Addr:    ":9090",
		Handler: &router{},
	}
	srv.ListenAndServe()
}

type router struct{}

var (
	healthzCode = 20
	VERSION     = "VERSION"
)

// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
// 3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
// 4. 当访问 localhost/healthz 时，应返回 20
func (*router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("client ip:%s path: %s\n", getRealIP(r), path)
	if path == "/healthz" {
		fmt.Printf("http response code:200\n")
		fmt.Fprintln(w, healthzCode)
	} else {
		// handle other request
		version := getOSEnv(VERSION, "")
		w.Header().Set(VERSION, version)
		fmt.Printf("http response code:200\n")
		fmt.Fprintln(w, "every thing's be ok")
	}
}

func getOSEnv(key, defaultValue string) string {
	result := os.Getenv(key)
	if result != "" {
		return result
	}
	return defaultValue
}

func getRealIP(r *http.Request) string {
	return r.RemoteAddr
}

type Server interface {
	Serve()
}
