package main

import (
	"flag"
	"lib"
	"log"
	"time"
)

// web服务默认监听端口
var port = 3500

func run() {
	var startserver bool
	flag.BoolVar(&startserver, "startserver", false, "Start server")
	flag.BoolVar(&startserver, "s", false, "Start server [shorted]")
	port1 := flag.Int("port", port, "Web server port >1024")

	var build bool
	flag.BoolVar(&build, "build", false, "Build")
	flag.BoolVar(&build, "b", false, "Build [shorted]")

	var publish bool
	flag.BoolVar(&publish, "publish", false, "Publish")
	flag.BoolVar(&publish, "p", false, "Publish [shorted]")
	var datetime bool
	flag.BoolVar(&datetime, "datetime", false, "Append datetime for Publish")
	flag.BoolVar(&datetime, "dt", false, "Append datetime for Publish [shorted]")

	flag.Parse()

	if startserver {
		if *port1 <= 1024 {
			log.Fatal("Web port must >1024")
		} else {
			lib.ServerHelper(*port1)
		}
	} else if build {
		lib.Build(true)
	} else if publish {
		if !lib.ReadConfig() {
			return
		}
		if datetime {
			lib.Config.Publish.Dir += time.Now().Format("20060102150405")
		}
		lib.Publish()
	} else {
		flag.PrintDefaults()
	}
}

func main() {
	run()
}
