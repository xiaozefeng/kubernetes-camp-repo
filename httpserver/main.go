package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/fsnotify/fsnotify"
)

// 有两种方式实现http server
// 1. 通过 golang 自带的 net/http server 实现， net/http 包已经实现了 http包的解析
// 2. 通过 golang 自带的 net/tcp 去实现， 这种方式需要自己解析http请求，自己拼装http响应

var cfg string

func init() {
	flag.StringVar(&cfg, "c", "", "config path")
	flag.Parse()
}

func initConfig(cfg string) error {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("conf")
	}
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("config file changed: %s ", e.Name)
		level, err := logrus.ParseLevel(viper.GetString("log.level"))
		if err != nil {
			logrus.Warn("日志配置出错")
			return
		}
		logrus.SetLevel(level)
		logrus.Infof("已经将日志等级改为 %v", viper.GetString("log.level"))
	})
	viper.WatchConfig()
	return nil
}

func initLog() error {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	return nil
}

func main() {
	err := initConfig(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	err = initLog()
	if err != nil {
		logrus.Fatal(err)
	}

	srv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: &Router{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logrus.Infof("strting http srever at address:%s \n", viper.GetString("addr"))
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-errCtx.Done()
		logrus.Infof("stopping http server\n")
		return srv.Shutdown(errCtx)
	})

	g.Go(func() error {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-done:
			logrus.Infof("sig: %v\n", sig)
			cancel()
		case <-errCtx.Done():
			return errCtx.Err()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logrus.Infof("err: %v\n", err)
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
		logrus.Infof("http response code:200\n")
		fmt.Fprintln(w, HEALTH_CODE)
	} else if path == "/debug" {
		logrus.Debug("this is debug request")
		fmt.Fprintln(w, HEALTH_CODE)
	} else {
		// handle other request
		version := getOSEnv(VERSION, "")
		w.Header().Set(VERSION, version)
		logrus.Infof("http response code:200\n")
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
