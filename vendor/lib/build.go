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
func buildCommand(isPublish bool) {
	cmdStr := "tsc"
	var filesLen = len(Config.Files)
	// 生成临时文件
	var tempFile = "ts-tool_temp.txt"
	if filesLen > 0 {
		var list []string
		for i := 0; i < len(Config.Files); i++ {
			item := Config.Files[i]
			item = strings.Replace(item, "ts", "js", 0)
			list = append(list, "src/"+item)
		}
		var str = strings.Join(list, "\r\n")
		ioutil.WriteFile(tempFile, ([]byte)(str), os.ModeAppend)
		cmdStr += " @" + tempFile
	}
	cmdStr += " --target " + Config.Target
	if isPublish {
		cmdStr += " --outFile " + SingleJs
	} else {
		cmdStr += " --outDir " + Config.OutJsDir
	}
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
	if filesLen > 0 {
		// 删除临时文件
		os.Remove(tempFile)
	}
}

// 生成html
func buildHtmls(isPublish bool) {
	var moduleJs string
	if len(Config.Modules) > 0 {
		for _, m := range Config.Modules {
			if isPublish {
				// 发布版的要放模块到 js/modules/ 下
				str := Config.OutJsDir + "/modules/" + m + ".min.js"
				moduleJs += "    <script src=\"" + str + "\"></script>\n"
			} else {
				str := Config.ModulesDir + "/" + m + "/" + m + ".js"
				moduleJs += "    <script src=\"" + str + "\"></script>\n"
			}
		}
	}
	var js string
	if isPublish {
		str := Config.OutJsDir + "/" + Config.Publish.MinJs
		js += "    <script src=\"" + str + "\"></script>\n"
	} else {
		for i := 0; i < len(Config.Files); i++ {
			str := Config.Files[i]
			str = strings.Replace(str, ".ts", ".js", -1)
			str = Config.OutJsDir + "/" + str
			js += "    <script src=\"" + str + "\"></script>\n"
		}
	}
	for i := 0; i < len(Config.Htmls); i++ {
		var html = Config.Htmls[i]
		if isPublish {
			html = Config.Publish.Dir + "/" + html
		}
		handleHTML(html, moduleJs, js)
	}
}

// 处理单个html文件
// 模块区块<!--modules_files_start--><!--modules_files_end-->
// js区块<!--game_files_start--><!--game_files_end-->
func handleHTML(filename string, moduleJs string, js string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("[Build html error]  %v\n", err)
		return
	}
	buffer := bufio.NewReader(f)
	var data string
	// 开始的类别 0没有开始 1模块开始 2js开始
	var typeflag = 0
	for {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			data += line
			break
		}
		if err == nil {
			if strings.Contains(line, "<!--modules_files_start-->") {
				data += line
				data += moduleJs
				typeflag = 1
			} else if strings.Contains(line, "<!--game_files_start-->") {
				data += line
				data += js
				typeflag = 2
			} else if strings.Contains(line, "<!--modules_files_end-->") || strings.Contains(line, "<!--game_files_end-->") {
				typeflag = 0
			}
			if typeflag == 0 {
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
