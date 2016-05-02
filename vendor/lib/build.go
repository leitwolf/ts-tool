package lib

//
// 构建生成
//
import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// SingleJs 只导出到一个单个文件
var SingleJs = "ts-tool_single.js"

// 执行生成命令
func buildCommand(isBuildSingle bool) {
	// 生成临时文件
	var tempFile = "ts-tool_temp.txt"
	var list []string
	for i := 0; i < len(Config.Files); i++ {
		item := Config.Files[i]
		item = strings.Replace(item, "ts", "js", 0)
		list = append(list, "src/"+item)
	}
	var str = strings.Join(list, "\r\n")
	ioutil.WriteFile(tempFile, ([]byte)(str), os.ModeAppend)
	cmdStr := "tsc @" + tempFile + " --target es5"
	if isBuildSingle {
		cmdStr += " --outFile " + SingleJs
	} else {
		cmdStr += " --outDir " + Config.JsDir
	}
	cmdList := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdList[0], cmdList[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("err", err)
	}
	str = string(out)
	if str != "" {
		log.Println(str)
	}
	// 删除临时文件
	os.Remove(tempFile)
}

// 生成html
func buildHtmls(isBuildSingle bool) {
	var js string
	if isBuildSingle {
		str := Config.JsDir + "/" + Config.Publish.MinJs
		js += "    <script src=\"" + str + "\"></script>\n"
	} else {
		for i := 0; i < len(Config.Files); i++ {
			str := Config.Files[i]
			str = strings.Replace(str, ".ts", ".js", -1)
			str = Config.JsDir + "/" + str
			js += "    <script src=\"" + str + "\"></script>\n"
		}
	}
	for i := 0; i < len(Config.Htmls); i++ {
		var html = Config.Htmls[i]
		if isBuildSingle {
			html = Config.Publish.Dir + "/" + html
		}
		handleHTML(html, js)
	}
}

// 处理单个html文件
func handleHTML(filename string, js string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("[Build html]  %v\n", err)
		return
	}
	buffer := bufio.NewReader(f)
	var data string
	// 是否要包含进来，在<!-- start和<!-- end之间的不包含
	var valid = true
	for {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			data += line
			break
		}
		if err == nil {
			if strings.Contains(line, "<!-- start") || strings.Contains(line, "<!--start") {
				data += line
				data += js
				valid = false
			} else if strings.Contains(line, "<!-- end") || strings.Contains(line, "<!--end") {
				valid = true
			}
			if valid {
				data += line
			}
		}
	}
	ioutil.WriteFile(filename, ([]byte)(data), os.ModeAppend)
}

// Build 编译生成
func Build() {
	if !ReadConfig() {
		return
	}
	buildCommand(false)
	buildHtmls(false)
}
