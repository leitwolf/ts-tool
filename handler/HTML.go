package handler

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"mo/conf"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
处理html
*/

// GetMinJs 获取min.js版本路径
func GetMinJs(jsPath string) string {
	l := len(jsPath)
	if !strings.HasSuffix(jsPath, ".min.js") && strings.HasSuffix(jsPath, ".js") {
		return jsPath[0:l-3] + ".min.js"
	}
	return jsPath
}

// CheckFileExists 判断文件是否存在
func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}

// HTML 处理html
type HTML struct {
}

// Build 生成html
func (h *HTML) Build(isPublish bool) {
	var libJs string
	var appJs string
	// 加上版本号（当前时间戳）
	verStr := strconv.FormatInt(time.Now().Unix(), 10)
	// 库文件列表
	for i := 0; i < len(conf.Conf.Libs); i++ {
		str := conf.Conf.Libs[i]
		// 真实路径
		if isPublish {
			// 优先min文件
			path1 := conf.Conf.Publish.Dir + "/" + GetMinJs(str)
			if CheckFileExists(path1) {
				str = GetMinJs(str)
			}
		}
		libJs += "    <script src=\"" + str + "?v=" + verStr + "\"></script>\n"
	}
	str := conf.Conf.OutJs
	if isPublish {
		str = GetMinJs(str)
	}
	appJs += "    <script src=\"" + str + "?v=" + verStr + "\"></script>\n"

	for i := 0; i < len(conf.Conf.HTML.List); i++ {
		var html = conf.Conf.HTML.List[i]
		if isPublish {
			html = conf.Conf.Publish.Dir + "/" + html
		}
		h.handleOne(html, libJs, appJs)
	}
}

// 处理单个html文件
// libJs 库文件区块
// appJs 程序区块
func (h *HTML) handleOne(filename string, libJs string, appJs string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("[Build html error]  %v\n", err)
		return
	}
	buffer := bufio.NewReader(f)
	var data string
	// 开始的类别 0 没有开始 1 libJs开始 2 appJs开始
	var typeflag = 0
	for {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			data += line
			break
		}
		if err == nil {
			if strings.Contains(line, conf.Conf.HTML.LibStartFlag) {
				data += line
				data += libJs
				typeflag = 1
			} else if strings.Contains(line, conf.Conf.HTML.AppStartFlag) {
				data += line
				data += appJs
				typeflag = 2
			} else if strings.Contains(line, conf.Conf.HTML.LibEndFlag) || strings.Contains(line, conf.Conf.HTML.AppEndFlag) {
				typeflag = 0
			}
			if typeflag == 0 {
				data += line
			}
		}
	}
	ioutil.WriteFile(filename, ([]byte)(data), os.ModeAppend)
}
