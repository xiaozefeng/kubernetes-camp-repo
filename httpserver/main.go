package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// 有两种方式实现http server
// 1. 通过 golang 自带的 net/http server 实现， net/http 包已经实现了 http包的解析
// 2. 通过 golang 自带的 net/tcp 去实现， 这种方式需要自己解析http请求，自己拼装http响应
var addr string

func main() {
	flag.StringVar(&addr, "a", ":9090", "http server addr")
	flag.Parse()

	srv := &http.Server{
		Addr:    addr,
		Handler: &Router{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Printf("strting http srever at address:%s \n", addr)
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-errCtx.Done()
		log.Printf("stopping http server\n")
		return srv.Shutdown(errCtx)
	})

	g.Go(func() error {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-done:
			log.Printf("sig: %v\n", sig)
			cancel()
		case <-errCtx.Done():
			return errCtx.Err()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("err: %v\n", err)
	}
}


const (
	HEALTH_CODE = 200
	VERSION     = "VERSION"
)

// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
// 3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
// 4. 当访问 localhost/healthz 时，应返回 20
type Router struct{}

func (*Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("client ip:%s path: %s\n", getRealIP(r), path)
	if path == "/healthz" {
		fmt.Printf("http response code:200\n")
		fmt.Fprintln(w, HEALTH_CODE)
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

