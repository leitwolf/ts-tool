# ts-tool
Typescript项目开发工具

# go版本
1.6以上

# 依赖
* github.com/fsnotify/fsnotify
* github.com/tdewolff/minify

# 特性
* 运行web服务器并监听资源文件及源码的更改，以及时的编译更新整个项目
* 构建项目
* 发布项目

# 参数

```shell
Usage of ts-tool:
  -b    Build [shorted]
  -build
        Build
  -datetime
        Append datetime for Publish
  -dt
        Append datetime for Publish [shorted]
  -p    Publish [shorted]
  -port int
        Web server port >1024 (default 3500)
  -publish
        Publish
  -s    Start server [shorted]
  -startserver
        Start server
```
运行web服务器并监听：

```shell
ts-tool -s -port 3000
```
构建：

```shell
ts-tool -b
```
发布：

```shell
ts-tool -p
```
> 其中port是跟startserver一起的
> datetime跟publish一起同，意思为是否在发布的目录名后面加上日期时间

# 配置文件

```json
{
    "target": "es5",
    "srcDir":"src",
    "outJs": "js/main.js",
    "res": {
        "path": "src/R.ts",
        "dir": "images"
    },
    "libs": [
        "js/libs/a.js",
        "js/libs/b.js"
    ],
    "html": {
        "libStartFlag": "<!--libs_start-->",
        "libEndFlag": "<!--libs_end-->",
        "appStartFlag": "<!--app_start-->",
        "appEndFlag": "<!--app_end-->",
        "list": [
            "index.html"
        ]
    },
    "publish": {
        "dir": "build",
        "copyList": [
            "js/libs",
            "images"
        ]
    }
}
```
* `target` 编译参数，默认 `es5` 
* `srcDir` 源码所在目录，默认 `src` 
* `outJs` 编译输出文件，默认 `js/main.js` 
* `res` 资源处理相关参数,此功能是把资源文件列到src/R.ts里，以便程序调用，默认空
* ----`path` 生成的ts文件路径，默认`src/R.ts`
* ----`dir` 需要处理的目录
* `libs` 引用第三方库列表
* `html` 处理html文件相关参数
* ----`libStartFlag` 第三方库起始标记，默认`<!--libs_start-->`
* ----`libEndFlag` 第三方库结束标记，默认`<!--libs_end-->`
* ----`appStartFlag` 生成js文件起始标记，默认`<!--app_start-->`
* ----`appEndFlag` 生成js文件结束标记，默认`<!--app_end-->`
* ----`list` 要处理的html文件列表，默认`["index.html"]`
* `publish` 发布相关参数，会复制第三方库的min文件（没有则复制原始的）
* ----`dir` 发布到的目录，默认`build`
* ----`copyList` 发布时直接拷贝的文件（或文件夹）列表，默认空

