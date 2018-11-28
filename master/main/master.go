package main

import (
	"runtime"
	"github.com/ywandy/crontab-go/master"
	"fmt"
	"flag"
)

var (
	confFile string //解析的配置文件路径
)

//命令行参数解析
func initArgs() {
	flag.StringVar()
}

//初始化一些环境
func initEnv() {
	//线程和cpu数量相同
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	//初始化线程
	initEnv()
	//加载配置

	//启动restful API HTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
ERR:
	fmt.Println(err)
}
