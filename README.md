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

# 配置文件

```json
{
    "target": "es5",
    "jsDir": "js",
    "resourceDir": "images",
    "publish": {
        "dir": "build",
        "minJs": "main.min.js",
        "copyFiles": [
            "js/lib",
            "images"
        ]
    },
    "htmls": [
        "index.html"
    ],
    "files": [
        "a.ts",
        "b.ts"
    ]
}
```
* `target` 编译参数，默认 `es5` 
* `jsDir` 编译ts文件到目录，默认 `js` 
* `resourceDir` 需要处理的资源目录，此功能是把资源文件列到src/Res.ts里，以便程序调用，默认空
* `publish` 发布相关参数
* ----`dir` 发布到的目录，默认`build`
* ----`minJs` 压缩成单一js文件的名称，默认`main.min.js`
* ----`copyFiles` 发布时直接拷贝的文件，默认空
* `htmls` 构建时需要更改的html主文件，在html的内容`<!-- start`和`<!-- end`之间加入编译好的js文件，默认`index.html`
* `files` 需要编译的ts文件列表，在`src`目录中，要注意文件顺序，默认空

