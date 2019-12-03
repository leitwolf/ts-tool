package common

import (
	"flag"
	"log"
	"mo/conf"
	"mo/handler"
	"time"
)

/*
命令行参数
*/

// Cmd 命令行
type Cmd struct {
	startserver bool
	port        int
	build       bool
	publish     bool
	datetime    bool
}

// 初始化
func (c *Cmd) init() {
	flag.BoolVar(&c.startserver, "startserver", false, "Start server")
	flag.BoolVar(&c.startserver, "s", false, "Start server [shorted]")
	flag.IntVar(&c.port, "port", Port, "Web server port >1024")

	flag.BoolVar(&c.build, "build", false, "Build")
	flag.BoolVar(&c.build, "b", false, "Build [shorted]")

	flag.BoolVar(&c.publish, "publish", false, "Publish")
	flag.BoolVar(&c.publish, "p", false, "Publish [shorted]")
	flag.BoolVar(&c.datetime, "datetime", false, "Append datetime for Publish")
	flag.BoolVar(&c.datetime, "dt", false, "Append datetime for Publish [shorted]")
}

// 解析
func (c *Cmd) parse() {
	flag.Parse()

	if c.startserver {
		if c.port <= 1024 {
			log.Fatal("Web port must >1024")
		} else {
			ServerHelper(c.port)
		}
	} else if c.build {
		handler.GBuild.BuildAll()
	} else if c.publish {
		if !conf.Conf.ReadConfig() {
			return
		}
		if c.datetime {
			conf.Conf.Publish.Dir += time.Now().Format("20060102150405")
		}
		handler.GPublish.Publish()
	} else {
		flag.PrintDefaults()
	}
}

// StartCmd 开始解析命令行
func StartCmd() {
	cmd := &Cmd{}
	cmd.init()
	cmd.parse()
}
