package lib

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

//
// GetIntranetIP 获取内网IP地址
//
func GetIntranetIP() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

//
// StartServer 启动服务器
//
func StartServer(port int) {
	// 内网地址
	addr := ":" + strconv.Itoa(port)
	internal := GetIntranetIP()
	if internal == "" {
		log.Println("Can not get internal ip.")
	} else {
		url := "http://" + internal + addr
		log.Println("Listening: " + url)
	}
	http.Handle("/", http.FileServer(http.Dir("./")))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// ServerHelper 启动web服务及watch服务
func ServerHelper(port int) {
	done := make(chan bool)
	Watch()
	defer watcher.Close()
	HandleImages()
	StartServer(port)
	// 阻塞退出（好像没什么用，web已经阻塞了）
	<-done
}
