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
class Res {
    
{{content}} 
    constructor() {
    }
}`

// 每一行的模板
const lineTmpl = "    public static {{name}} : string = \"{{path}}\";\r\n"

// 保存的位置
const savePath = "src/Res.ts"

// 图片目录
var imagePath = "images/"

//
// HandleImages 处理资源
// images/select/box2.png 名字为 select_box2
//
func HandleImages() {
	path := imagePath
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
		// 只有图片才处理
		// ext := path1[dotIndex:]
		// if ext != ".png" && ext != ".jpg" {
		// 	return nil
		// }
		name := path1[l:dotIndex]
		name = strings.Replace(name, "/", "_", -1)
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
