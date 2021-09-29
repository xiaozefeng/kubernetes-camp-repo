package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
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

	if req.Path == "/healthz" {
		write(conn, req, []byte("200"))
	} else {
		write(conn, req, []byte(req.Path))
	}
}

func getResponseHeader(header map[string]string) string {
	var buf bytes.Buffer
	for k, v := range header {
		if k == "Accept" || k == "User-Agent"{
			continue
		}
		buf.WriteString(k)
		buf.WriteByte(':')
		buf.WriteByte(' ')
		buf.WriteString(v)
		buf.WriteByte('\n')
	}
	res:=buf.String()
	return res[:len(res)-1]
}

func write(conn net.Conn, req *HTTPRequest, v []byte) (int, error) {
	req.Header["Date"] = time.Now().String()
	req.Header["Content-Length"] = fmt.Sprintf("%d", len(v))
	req.Header["Content-Type"] = "text/plain;charset=UTF-8"
	req.Header[VERSION] = getOSEnv(VERSION, "")

	var header = getResponseHeader(req.Header)
	data := fmt.Sprintf("HTTP/1.1 200 OK\n%s\n\n%s\n",
		header,
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
	// Body        io.Reader
	ReomoteAddr string
}
