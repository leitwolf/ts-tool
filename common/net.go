package common

import (
	"mo/conf"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	//Port web服务默认监听端口
	Port = 3500
)

// GetIntranetIP 获取内网IP地址
func GetIntranetIP() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

// StartServer 启动服务器
func StartServer(port int) {
	// 内网地址
	addr := ":" + strconv.Itoa(port)
	host := GetIntranetIP()
	if host == "" {
		host = "127.0.0.1"
	}
	url := "http://" + host + addr
	log.Println("Listening: " + url)
	// 静态文件
	http.Handle("/", http.FileServer(http.Dir("./")))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// ServerHelper 启动web服务及watch服务
func ServerHelper(port int) {
	// 加载配置文件
	conf.Conf.ReadConfig()
	done := make(chan bool)
	GWatch.Watch()
	defer GWatch.watcher.Close()
	StartServer(port)
	// 阻塞退出（好像没什么用，web已经阻塞了）
	<-done
}
