package common

/*
监视文件变化及处理
*/

import (
	"conf"
	"handler"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch 监视器
type Watch struct {
	// 监听器实例
	watcher *fsnotify.Watcher
	// 是否需要处理资源
	needHandleRes bool
	// 是否需要构建
	needBuild bool
	// 已构建次数
	buildCount int
	// 计时器，每隔一分钟，如果没有动静则会构建一次
	timer *time.Timer
}

// 添加监听
func (w *Watch) addWatch(path string) {
	log.Println("Watching path: ", path)
	err := w.watcher.Add(path)
	if err != nil {
		log.Printf("Watch "+path+" error: %v\n", err)
	}
}

// 删除监听
func (w *Watch) removeWatch(path string) {
	log.Println("Remove watch path: ", path)
	err := w.watcher.Remove(path)
	if err != nil {
		// log.Printf("Remove watch "+path+" error: %v\n", err)
	}
}

// 资源处理器
func (w *Watch) startResHandler() {
	go func() {
		for {
			if w.needHandleRes {
				w.needHandleRes = false
				log.Println("Handling res...")
				handler.GRes.Handle()
				log.Println("Handle res done")
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

// 构建处理器
func (w *Watch) startBuildHandler() {
	go func() {
		for {
			if w.needBuild {
				w.needBuild = false
				w.buildCount++
				log.Println("Building...", w.buildCount)
				handler.GBuild.Build()
				log.Println("Build done")
				w.startTick()
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

// 60秒执行一次
func (w *Watch) startTick() {
	if w.timer != nil {
		w.timer.Stop()
	}
	w.timer = time.NewTimer(time.Second * 60)
	go func() {
		<-w.timer.C
		w.needBuild = true
		log.Println("Building in period")
	}()
}

// 处理新建监听
func (w *Watch) watchCreate(path string) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}
	srcPath := conf.Conf.SrcDir
	// resPath := Config.ResourceDir
	if file.IsDir() {
		// 建立目录
		w.addWatch(path)
	} else if strings.HasPrefix(path, srcPath) {
		// 在源码目录中建立文件
		w.needBuild = true
	}
}

// 处理删除监听
func (w *Watch) watchRemove(path string) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}
	srcPath := conf.Conf.SrcDir
	if file.IsDir() {
		// 删除目录
		w.removeWatch(path)
	} else if strings.HasPrefix(path, srcPath) {
		// 在源码目录中删除文件
		w.needBuild = true
	}
}

// 处理重命名监听
func (w *Watch) watchRename(path string) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}
	srcPath := conf.Conf.SrcDir
	if file.IsDir() {
		// 修改目录名称
		w.addWatch(path)
	} else if strings.HasPrefix(path, srcPath) {
		// 在源码目录中修改文件名
		w.needBuild = true
	}
}

// 处理修改监听
func (w *Watch) watchWrite(path string) {
	srcPath := conf.Conf.SrcDir
	if strings.HasPrefix(path, srcPath) {
		// 源码修改
		w.needBuild = true
	}
}

// 处理监听到的变化
func (w *Watch) handleWatch(event fsnotify.Event) {
	path := strings.Replace(event.Name, "\\", "/", -1)
	if strings.Contains(path, conf.ConfigFilename) {
		// 配置文件改变
		w.needBuild = true
	} else if w.checkResChanged(path) {
		// 素材改变
		w.needHandleRes = true
	} else if event.Op&fsnotify.Create == fsnotify.Create {
		w.watchCreate(path)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		w.watchRemove(path)
	} else if event.Op&fsnotify.Rename == fsnotify.Rename {
		w.watchRename(path)
	} else if event.Op&fsnotify.Write == fsnotify.Write {
		w.watchWrite(path)
	}
}

// 检测是否需要更新资源
func (w *Watch) checkResChanged(filename string) bool {
	resDirs := conf.Conf.GetResDirs()
	for i := 0; i < len(resDirs); i++ {
		if strings.HasPrefix(filename, resDirs[i]) {
			return true
		}
	}
	return false
}

// Watch 开始监听
func (w *Watch) Watch() {
	nw, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Watch new error: %v\n", err)
		return
	}
	w.watcher = nw
	// defer watcher.Close()

	// done := make(chan bool)
	// 开启一个 goroutine 来处理事件
	go func() {
		for {
			select {
			case event := <-w.watcher.Events:
				if event.Name != "" {
					// log.Println("Watch event:", event)
					w.handleWatch(event)
				}
			case err := <-w.watcher.Errors:
				if err != nil {
					log.Printf("Watch error: %v\n", err)
				}
			}
		}
	}()

	// 监听的目录列表，包含子目录
	list := []string{conf.Conf.SrcDir}
	resDirs := conf.Conf.GetResDirs()
	for i := 0; i < len(resDirs); i++ {
		list = append(list, resDirs[i])
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
	allList = append(allList, conf.ConfigFilename)
	for i := 0; i < len(allList); i++ {
		item := allList[i]
		w.addWatch(item)
	}
	// 开启各个处理器
	w.startResHandler()
	w.startBuildHandler()
	// <-done
}

// GWatch Watch 单例
var GWatch = &Watch{needBuild: true, needHandleRes: true, buildCount: 0}
