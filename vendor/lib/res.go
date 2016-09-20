package lib

//
// 处理资源文件，生成一个Res.ts的文件到src目录
//
import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 资源模板
const template = `/**
 * DO NOT edit
 */
class R {
    
{{content}} 
    constructor() {
    }
}`

// 每一行的模板
const lineTmpl = "    public static {{name}} : string = \"{{path}}\";\r\n"

// 保存的位置
const savePath = "src/R.ts"

//
// HandleImages 处理资源
// images/select/box2.png 名字为 select_box2
//
func HandleImages() {
	if Config.ResourceDir == "" {
		return
	}
	path := Config.ResourceDir
	content := ""
	l := len(path)
	err := filepath.Walk(path, func(path1 string, f os.FileInfo, err1 error) error {
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
	content = strings.Replace(template, "{{content}}", content, -1)
	ioutil.WriteFile(savePath, ([]byte)(content), os.ModeAppend)
}
