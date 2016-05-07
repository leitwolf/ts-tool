package lib

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// 监听器实例
var watcher *fsnotify.Watcher

// 添加监听
func addWatch(path string) {
	log.Println("Watching path: ", path)
	err := watcher.Add(path)
	if err != nil {
		log.Printf("Watch "+path+" error: %v\n", err)
	}
}

// 删除监听
func removeWatch(path string) {
	log.Println("Remove watch path: ", path)
	err := watcher.Remove(path)
	if err != nil {
		// log.Printf("Remove watch "+path+" error: %v\n", err)
	}
}

// 是否需要处理资源
var needHandleRes = false

// 资源处理器
func startResHandler() {
	go func() {
		for {
			if needHandleRes {
				needHandleRes = false
				log.Println("Handling res...")
				HandleImages()
				log.Println("Handle res done")
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

// 是否需要构建
var needBuild = false

// 已构建次数
var buildCount = 0

// 构建处理器
func startBuildHandler() {
	go func() {
		for {
			if needBuild {
				needBuild = false
				buildCount++
				log.Println("Building...", buildCount)
				Build(false)
				log.Println("Build done")
				startTick()
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

// 计时器，每隔一分钟，如果没有动静则会构建一次
var timer *time.Timer

func startTick() {
	if timer != nil {
		timer.Stop()
	}
	timer = time.NewTimer(time.Second * 60)
	go func() {
		<-timer.C
		needBuild = true
		log.Println("Building in period")
	}()
}

// 处理新建监听
func watchCreate(path string) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}
	srcPath := "src"
	resPath := Config.ResourceDir
	if file.IsDir() {
		// 建立目录
		addWatch(path)
	} else if resPath != "" && strings.HasPrefix(path, resPath) {
		// 在资源目录中建立文件
		needHandleRes = true
	} else if strings.HasPrefix(path, srcPath) {
		// 在源码目录中建立文件
		needBuild = true
	}
}

// 处理删除监听
func watchRemove(path string) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}
	srcPath := "src"
	resPath := Config.ResourceDir
	if file.IsDir() {
		// 删除目录
		removeWatch(path)
	} else if resPath != "" && strings.HasPrefix(path, resPath) {
		// 在资源目录中删除文件
		needHandleRes = true
	} else if strings.HasPrefix(path, srcPath) {
		// 在源码目录中删除文件
		needBuild = true
	}
}

// 处理重命名监听
func watchRename(path string) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}
	srcPath := "src"
	resPath := Config.ResourceDir
	if file.IsDir() {
		// 修改目录名称
		addWatch(path)
	} else if resPath != "" && strings.HasPrefix(path, resPath) {
		// 在资源目录中修改文件名
		needHandleRes = true
	} else if strings.HasPrefix(path, srcPath) {
		// 在源码目录中修改文件名
		needBuild = true
	}
}

// 处理修改监听
func watchWrite(path string) {
	srcPath := "src"
	if strings.HasPrefix(path, srcPath) {
		// 源码修改
		needBuild = true
	}
}

// 处理监听
func handleWatch(event fsnotify.Event) {
	path := strings.Replace(event.Name, "\\", "/", -1)
	if strings.Contains(path, configFilename) {
		// 配置文件改变
		needBuild = true
	} else if Config.ModulesDir != "" && strings.Contains(path, Config.ModulesDir) {
		// 模块改变
		needBuild = true
	} else if event.Op&fsnotify.Create == fsnotify.Create {
		watchCreate(path)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		watchRemove(path)
	} else if event.Op&fsnotify.Rename == fsnotify.Rename {
		watchRename(path)
	} else if event.Op&fsnotify.Write == fsnotify.Write {
		watchWrite(path)
	}
}

// Watch 开始监听
func Watch() {
	// 先读取配置文件
	ReadConfig()

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Watch new error: %v\n", err)
		return
	}
	watcher = w
	// defer watcher.Close()

	// done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Name != "" {
					log.Println("Watch event:", event)
					handleWatch(event)
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Printf("Watch error: %v\n", err)
				}
			}
		}
	}()

	// 监听的目录列表，包含子目录
	list := []string{"src"}
	if Config.ResourceDir != "" {
		list = append(list, Config.ResourceDir)
	}
	if Config.ModulesDir != "" {
		list = append(list, Config.ModulesDir)
	}
	// log.Printf("%v\n", list)
	var allList []string
	for i := 0; i < len(list); i++ {
		item := list[i]
		filepath.Walk(item, func(path1 string, f os.FileInfo, err1 error) error {
			if f == nil {
				return err1
			}
			if f.IsDir() {
				allList = append(allList, path1)
				return nil
			}
			return nil
		})
	}
	// 添加配置文件监听
	allList = append(allList, configFilename)
	for i := 0; i < len(allList); i++ {
		item := allList[i]
		addWatch(item)
	}
	// 开启各个处理器
	startBuildHandler()
	startResHandler()

	// <-done
}
