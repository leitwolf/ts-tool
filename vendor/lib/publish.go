package lib

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

// 复制文件
func copyFiles() {
	list := Config.Publish.CopyFiles
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
				destPath := Config.Publish.Dir + "/" + srcPath
				copyFile(srcPath, destPath)
				return nil
			})
			if err != nil {
				log.Printf("[copyFile] error %v\n", err)
			}
		} else {
			srcPath := path
			destPath := Config.Publish.Dir + "/" + file.Name()
			copyFile(srcPath, destPath)
		}
	}
	// 复制html
	list = Config.Htmls
	for i := 0; i < len(list); i++ {
		srcPath := list[i]
		destPath := Config.Publish.Dir + "/" + srcPath
		copyFile(srcPath, destPath)
	}
	// 复制模块
	if len(Config.Modules) > 0 {
		for _, m := range Config.Modules {
			srcPath := Config.ModulesDir + "/" + m + "/" + m + ".min.js"
			destPath := Config.Publish.Dir + "/" + Config.OutJsDir + "/modules/" + m + ".min.js"
			copyFile(srcPath, destPath)
			// 有可能有 web.min.js
			webPath := Config.ModulesDir + "/" + m + "/" + m + ".web.min.js"
			_, err := os.Stat(webPath)
			if err == nil {
				destPath2 := Config.Publish.Dir + "/" + Config.OutJsDir + "/modules/" + m + ".web.min.js"
				copyFile(webPath, destPath2)
			}
		}
	}
}

// 复制文件
func copyFile(src string, dest string) {
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
func minifyJs(inputFile string, outFile string) bool {
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	bytes, _ := ioutil.ReadFile(inputFile)
	b, err := m.Bytes("text/javascript", bytes)
	if err != nil {
		log.Printf("[minify error]%v", err)
		return false
	}
	ioutil.WriteFile(outFile, b, os.ModeAppend)
	os.Remove(inputFile)
	return true
}

// Publish 发布
// @param datetime 是否在目录后面加上当前时间
func Publish() {
	// 先删除文件
	os.RemoveAll(Config.Publish.Dir)
	copyFiles()
	buildCommand(true)
	buildHtmls(true)
	var destJs = Config.Publish.Dir + "/" + Config.OutJsDir + "/" + Config.Publish.MinJs
	minifyJs(SingleJs, destJs)
}
