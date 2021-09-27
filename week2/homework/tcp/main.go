package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var addr string

func main() {
	flag.StringVar(&addr, "a", ":9090", "tcp server addr")
	flag.Parse()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("starting tcp server at addr: %s", addr)
	defer l.Close()

	for {
		connection, err := l.Accept()
		if err != nil {
			log.Printf("connection happens eror:%v\n", err)
		}

		go handleConnection(context.TODO(), connection)
	}

}

const (
	VERSION = "VERSION"
)

func getOSEnv(key, defaultValue string) string {
	result := os.Getenv(key)
	if result != "" {
		return result
	}
	return defaultValue
}
func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr()
	log.Printf("remoteAddr: %s\n", remoteAddr)
	req, err := parseHTTP(conn)
	if err != nil {
		log.Printf("parse http happens error:%v \n", err)
		return
	}
	log.Printf("req: %+v\n", req)
	version := getOSEnv(VERSION, "")
	log.Printf("version: %v\n", version)
	if req.Path == "/healthz" {
		write(conn, []byte("20"), version)
	} else {
		write(conn, []byte(req.Path), version)
	}
}

func write(conn net.Conn, v []byte, version string) (int, error) {
	data := fmt.Sprintf("HTTP/1.1 200 OK\nDate: %v\nContent-Length:%d\nContent-Type: %s\nVersion:%s\n\n%s\n",
		time.Now(),
		len(v),
		"text/plain;charset=UTF-8",
		version,
		v)
	return conn.Write([]byte(data))
}

func parseHTTP(conn net.Conn) (*HTTPRequest, error) {

	reader := bufio.NewReader(conn)
	var err error
	var lines []string
	for err == nil {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("read line happens error: %v\n", err)
			break
		}
		if line == "\r\n" {
			break
		}
		lines = append(lines, line)
	}

	for _, v := range lines {
		fmt.Printf("%v", v)
	}

	var req = &HTTPRequest{}

	if len(lines) > 0 {
		firstLine := lines[0]
		splited := strings.Split(firstLine, ` `)
		req.Method = splited[0]
		req.Path = splited[1]
		req.Protocal = splited[2]
	}
	lines = lines[1:]

	if len(lines) > 0 {
		var header = make(map[string]string)
		for _, line := range lines {
			splited := strings.Split(line, `:`)
			header[strings.TrimRight(splited[0], ` `)] = strings.TrimLeft(splited[1], ` `)
		}
		req.Header = header
	}
	req.ReomoteAddr = conn.RemoteAddr().String()

	return req, nil
}

type HTTPRequest struct {
	Method      string
	Path        string
	Protocal    string
	Header      map[string]string
	Body        io.Reader
	ReomoteAddr string
}
