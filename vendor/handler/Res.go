package handler

/*
资源处理相关
*/

import (
	"conf"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ResTemplate 资源模板
	ResTemplate = `/**
 * DO NOT edit
 */
class R {
    
{{content}} 
    constructor() {
    }
}`
	lineTmpl = "    public static {{name}} : string = \"{{path}}\";\r\n"
)

// Res 资源处理
type Res struct {
}

// Handle 处理资源
// images/select/box2.png 名字为 select_box2
func (r *Res) Handle() {
	if conf.Conf.Res.Dir == "" || conf.Conf.Res.Path == "" {
		return
	}
	resDir := conf.Conf.Res.Dir
	content := ""
	l := len(resDir)
	err := filepath.Walk(resDir, func(path1 string, f os.FileInfo, err1 error) error {
		if f == nil {
			return err1
		}
		if f.IsDir() {
			return nil
		}
		path1 = strings.Replace(path1, "\\", "/", -1)
		dotIndex := strings.LastIndex(path1, ".")
		name := path1[l:dotIndex]
		// 非图片加后缀
		ext := strings.ToLower(path1[dotIndex:])
		if ext != ".png" && ext != ".jpg" && ext != ".bmp" {
			ext = strings.Replace(ext, ".", "_", -1)
			name += ext
		}
		name = strings.Replace(name, "/", "_", -1)
		if name[0] == '_' {
			name = name[1:]
		}
		line := strings.Replace(lineTmpl, "{{name}}", name, -1)
		line = strings.Replace(line, "{{path}}", path1, -1)
		content += line
		return nil
	})
	if err != nil {
		log.Printf("filepath.Walk() return %v\n", err)
	}
	// 生成
	content = strings.Replace(ResTemplate, "{{content}}", content, -1)
	ioutil.WriteFile(conf.Conf.Res.Path, ([]byte)(content), os.ModeAppend)
}

// GRes Res 单例
var GRes = &Res{}
