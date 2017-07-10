package main

import (
	"common"
	"conf"
)

// 运行
func run() {
	conf.InitConfig()
	common.StartCmd()
}

func main() {
	run()
}
