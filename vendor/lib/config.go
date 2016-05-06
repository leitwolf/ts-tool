package lib

//
// 配置文件相关
//
import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// PublishConfig 发布相关的配置
type PublishConfig struct {
	Dir       string   `json:"dir"`
	MinJs     string   `json:"minJs"`
	CopyFiles []string `json:"copyFiles"`
}

// TstoolConfig 配置文件结构
// ModulesDir 为模块文件夹，每个模块由三个文件组成:模块名m，文件m.d.ts, m.js, m.min.js
type TstoolConfig struct {
	Target      string `json:"target"`
	OutJsDir    string `json:"outJsDir"`
	ResourceDir string `json:"resourceDir"`
	ModulesDir  string `json:"modulesDir"`
	Modules     []string
	Publish     PublishConfig `json:"publish"`
	Htmls       []string      `json:"htmls"`
	Files       []string      `json:"files"`
}

// Config 配置文件
var Config TstoolConfig

// 配置文件名
var configFilename = "ts-toolconfig.json"

// 默认配置
func intConfig() {
	publish := PublishConfig{
		Dir:   "build",
		MinJs: "main.min.js",
	}
	Config = TstoolConfig{
		Target:      "es5",
		OutJsDir:    "js",
		ResourceDir: "",
		ModulesDir:  "",
		Modules:     []string{},
		Publish:     publish,
		Htmls:       []string{"index.html"},
	}
}

// ReadConfig 读取配置文件
func ReadConfig() bool {
	bytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Printf("Read config file error:%v\n", err)
		return false
	}
	intConfig()
	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		log.Println("Parse config file error:", err)
		return false
	}
	// 处理模块
	if Config.ModulesDir != "" {
		files, err := ioutil.ReadDir(Config.ModulesDir)
		if err != nil {
			log.Println("[Read Module dir error]", err)
		}
		for _, fi := range files {
			if fi.IsDir() {
				Config.Modules = append(Config.Modules, fi.Name())
			}
		}
	}
	// log.Printf("%v", Config)
	return true
}
