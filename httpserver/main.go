package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/fsnotify/fsnotify"
	"httpserver/mertrics"
	_ "net/http/pprof"
)

// 有两种方式实现http server
// 1. 通过 golang 自带的 net/http server 实现， net/http 包已经实现了 http包的解析
// 2. 通过 golang 自带的 net/tcp 去实现， 这种方式需要自己解析http请求，自己拼装http响应

var cfg string

func init() {
	flag.StringVar(&cfg, "c", "conf/conf.yaml", "config path")
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

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/healthz", Healthz)
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: mux,
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

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	version := getOSEnv(VERSION, "")
	w.Header().Set(VERSION, version)
	timer := mertrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	logrus.Infof("http response code:200\n")
	_, _ = fmt.Fprintln(w, "every thing's be ok")

}
func Healthz(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("client ip:%s path: %s\n", getRealIP(r), path)
	logrus.Infof("http response code:200\n")
	_, _ = fmt.Fprintln(w, HEALTH_CODE)
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

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
