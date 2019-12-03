package handler

/*
资源处理相关
*/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mo/conf"
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

// SpriteSheet 结构
type SpriteSheet struct {
	// File 图片文件名
	File string `json:"file"`
	// Frames 帧列表
	Frames map[string]interface{} `json:"frames"`
}

// Res 资源处理
type Res struct {
}

// Handle 处理资源
// images/select/box2.png 名字为 select_box2
// 其它文件夹如 sounds/bg.mp3 则为 sounds_bg
func (r *Res) Handle() {
	resDirs := conf.Conf.GetResDirs()
	if len(resDirs) == 0 || conf.Conf.Res.Path == "" {
		return
	}
	content := ""
	for i := 0; i < len(resDirs); i++ {
		content += r.handleDir(resDirs[i])
	}
	// 生成
	content = strings.Replace(ResTemplate, "{{content}}", content, -1)
	ioutil.WriteFile(conf.Conf.Res.Path, ([]byte)(content), os.ModeAppend)
}

// 处理一个文件夹
func (r *Res) handleDir(resDir string) string {
	content := ""
	dirLen := 0
	if resDir == "images" {
		dirLen = len(resDir)
	}
	err := filepath.Walk(resDir, func(path1 string, f os.FileInfo, err1 error) error {
		if f == nil {
			return err1
		}
		if f.IsDir() {
			return nil
		}
		filename := path1
		newPath := strings.Replace(filename, "\\", "/", -1)
		dotIndex := strings.LastIndex(newPath, ".")
		name := newPath[dirLen:dotIndex]
		// 非图片加后缀
		ext := strings.ToLower(newPath[dotIndex:])
		if ext != ".png" && ext != ".jpg" && ext != ".bmp" {
			name += strings.Replace(ext, ".", "_", -1)
		}
		name = strings.Replace(name, "/", "_", -1)
		if name[0] == '_' {
			name = name[1:]
		}
		line := strings.Replace(lineTmpl, "{{name}}", name, -1)
		line = strings.Replace(line, "{{path}}", newPath, -1)
		content += line

		// sprite sheet 的情况
		// log.Println(ext + " " + filename)
		if ext == ".json" {
			content += r.handleSheet(filename)
		}
		return nil
	})
	if err != nil {
		log.Printf("filepath.Walk() return %v\n", err)
	}
	return content
}

// 处理sprite sheet文件，把里面的项目也要加到资源列表中
func (r *Res) handleSheet(filename string) string {
	// log.Println("Parsing sprite sheet: " + filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Read sheet file error:", err.Error())
		return ""
	}
	result := SpriteSheet{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println("Parse sprite sheet file error:", err.Error())
		return ""
	}
	content := ""
	for k, _ := range result.Frames {
		value := k
		line := strings.Replace(lineTmpl, "{{name}}", value, -1)
		line = strings.Replace(line, "{{path}}", value, -1)
		content += line
	}
	return content
}

// Test 测试
func (r *Res) Test() {
	content := r.handleSheet("game.json")
	println(content)
}

// GRes Res 单例
var GRes = &Res{}
