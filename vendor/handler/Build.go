package handler

import (
	"conf"
	"log"
	"os/exec"
	"strings"
)

/*
编译生成
*/

// Build 编译
type Build struct {
	html *HTML
}

// Build 编译生成
func (b *Build) Build() {
	b.Build2(false, false)
}

// BuildAll 编译生成所有，包括资源
func (b *Build) BuildAll() {
	b.Build2(false, true)
}

// Build2 执行具体生成命令
func (b *Build) Build2(isPublish bool, handleRes bool) {
	if !conf.Conf.ReadConfig() {
		return
	}
	if handleRes {
		GRes.Handle()
	}
	b.genCmd(conf.Conf.OutJs)
	b.html.Build(isPublish)
}

func (b *Build) genCmd(outJs string) {
	cmdStr := "tsc"
	cmdStr += " --target " + conf.Conf.Target
	cmdStr += " --outFile " + outJs
	cmdList := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdList[0], cmdList[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("err", err)
	}
	outStr := string(out)
	if outStr != "" {
		log.Println(outStr)
	}
}

// GBuild Build 单例
var GBuild = &Build{html: &HTML{}}
