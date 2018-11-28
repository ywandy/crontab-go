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
	flag.StringVar(&confFile, "config", "./master.json", "指定传入的配置文件")
	//解析
	flag.Parse()
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
	//初始化命令行参数
	initArgs()
	//初始化线程
	initEnv()
	//加载配置文件
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	//启动任务管理器(etcd)
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	//启动restful API HTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
ERR:
	fmt.Println(err)
}
