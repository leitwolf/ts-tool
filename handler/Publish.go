package handler

/*
发布模块
*/

import (
	"io"
	"io/ioutil"
	"log"
	"mo/conf"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	minify "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

// Publish 发布模块
type Publish struct {
}

// 复制文件
func (p *Publish) copyFiles() {
	list := conf.Conf.Publish.CopyList
	for i := 0; i < len(list); i++ {
		path := list[i]
		// 检测是否是文件
		file, err := os.Stat(path)
		if err != nil {
			continue
		}
		if file.IsDir() {
			err := filepath.Walk(path, func(path1 string, f os.FileInfo, err1 error) error {
				if f == nil {
					return err1
				}
				if f.IsDir() {
					return nil
				}
				srcPath := strings.Replace(path1, "\\", "/", -1)
				destPath := conf.Conf.Publish.Dir + "/" + srcPath
				p.copyFile(srcPath, destPath)
				return nil
			})
			if err != nil {
				log.Printf("[copyFile] error %v\n", err)
			}
		} else {
			srcPath := path
			destPath := conf.Conf.Publish.Dir + "/" + file.Name()
			p.copyFile(srcPath, destPath)
		}
	}
	// 复制html
	list = conf.Conf.HTML.List
	for i := 0; i < len(list); i++ {
		srcPath := list[i]
		destPath := conf.Conf.Publish.Dir + "/" + srcPath
		p.copyFile(srcPath, destPath)
	}
	// 复制库文件
	if len(conf.Conf.Libs) > 0 {
		for _, item := range conf.Conf.Libs {
			srcPath := GetMinJs(item)
			destPath := conf.Conf.Publish.Dir + "/" + srcPath
			// 检测是否有对应的min文件，有则复制，没有则压缩并移到目标目录
			if !CheckFileExists(srcPath) {
				success := p.minifyJs(item, destPath)
				if !success {
					p.copyFile(item, destPath)
				}
			} else {
				p.copyFile(srcPath, destPath)
			}
		}
	}
}

// 复制文件
func (p *Publish) copyFile(src string, dest string) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Printf("Copy file "+src+" error %v\n", err)
		return
	}
	defer srcFile.Close()

	// 要先建立目录
	os.MkdirAll(path.Dir(dest), 0777)
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("Create file "+dest+" error %v\n", err)
		return
	}
	defer destFile.Close()
	io.Copy(destFile, srcFile)
}

// 压缩js
func (p *Publish) minifyJs(inputFile string, outFile string) bool {
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	bytes, _ := ioutil.ReadFile(inputFile)
	b, err := m.Bytes("text/javascript", bytes)
	if err != nil {
		log.Printf("[minify error]%v", err)
		return false
	}
	os.MkdirAll(path.Dir(outFile), 0777)
	ioutil.WriteFile(outFile, b, os.ModeAppend)
	return true
}

// Publish 发布
func (p *Publish) Publish() {
	// 先删除文件
	os.RemoveAll(conf.Conf.Publish.Dir)
	p.copyFiles()
	GBuild.Build2(true, true)
	// 生成临时文件
	srcPath := strconv.FormatInt(time.Now().UnixNano(), 10) + ".js"
	// println(srcPath)
	GBuild.genCmd(srcPath)
	destJs := conf.Conf.Publish.Dir + "/" + GetMinJs(conf.Conf.OutJs)
	p.minifyJs(srcPath, destJs)
	// 删除临时文件
	os.Remove(srcPath)
}

// GPublish Publish 单例
var GPublish = &Publish{}
