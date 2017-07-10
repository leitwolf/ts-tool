package conf

/*
配置相关
*/

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

// ConfigFilename 配置文件名
var ConfigFilename = "ts-toolconfig.json"

// ResConfig 资源文件相关的配置
type ResConfig struct {
	// 生成文件路径
	Path string `json:"path"`
	// 资源文件夹，监听变化以便生成R.ts文件
	Dir string `json:"dir"`
}

// HTMLConfig html文件相关的配置
type HTMLConfig struct {
	// 库文件开始与结束的标记
	LibStartFlag string `json:"libStartFlag"`
	LibEndFlag   string `json:"libEndFlag"`
	// 程序编译文件开始与结束的标记
	AppStartFlag string `json:"appStartFlag"`
	AppEndFlag   string `json:"appEndFlag"`
	// html文件列表
	List []string `json:"list"`
}

// PublishConfig 发布相关的配置
type PublishConfig struct {
	// 发布目录
	Dir string `json:"dir"`
	// 需要复制的文件(文件夹)列表
	CopyList []string `json:"copyList"`
}

// Config 总配置文件
type Config struct {
	Target string `json:"target"`
	// 源码目录
	SrcDir string `json:"srcDir"`
	// 输出文件（只生成一个文件）
	OutJs string `json:"outJs"`
	// 资源文件夹，监听变化以便生成R.ts文件
	Res ResConfig `json:"res"`
	// 库文件列表
	Libs    []string      `json:"libs"`
	HTML    HTMLConfig    `json:"html"`
	Publish PublishConfig `json:"publish"`
}

// 初始化默认配置
func (c *Config) init() {
	c.Res = ResConfig{
		Path: "src/R.ts",
		Dir:  "",
	}
	c.HTML = HTMLConfig{
		LibStartFlag: "<!--libs_start-->",
		LibEndFlag:   "<!--libs_end-->",
		AppStartFlag: "<!--app_start-->",
		AppEndFlag:   "<!--app_end-->",
		List:         []string{"index.html"},
	}
	c.Publish = PublishConfig{
		Dir:      "build",
		CopyList: []string{},
	}
	c.Target = "es5"
	c.SrcDir = "src"
	c.OutJs = "js/main.js"
	c.Libs = make([]string, 0)
}

// ReadConfig 读取配置文件
func (c *Config) ReadConfig() bool {
	data, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		log.Println("Read config file error:", err.Error())
		return false
	}
	// 过滤注释
	data, err = c.filterComments(data)
	if err != nil {
		log.Println("filter config file comments error:", err.Error())
		return false
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Println("Parse config file error:", err.Error())
		return false
	}
	return true
}

// 去掉配置文件中的注释
func (c *Config) filterComments(data []byte) ([]byte, error) {
	// windows
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0)
	lines := bytes.Split(data, []byte("\n"))
	var data2 [][]byte
	for _, line := range lines {
		match, err := regexp.Match(`^\s*/{2}`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			data2 = append(data2, line)
		}
	}
	return bytes.Join(data2, []byte("\n")), nil
}

// Conf 配置文件
var Conf *Config

// InitConfig 初始化配置
func InitConfig() {
	Conf = &Config{}
	Conf.init()
	// Conf.ReadConfig()
}
