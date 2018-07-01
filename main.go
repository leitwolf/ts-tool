package main

import (
	"common"
	"conf"
)

// 运行
func run() {
	conf.InitConfig()
	common.StartCmd()
	// handler.GRes.Test()
}

func main() {
	run()
}
