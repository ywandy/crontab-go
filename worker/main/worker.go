package main

import (
	"runtime"
	"fmt"
	"flag"
	"github.com/ywandy/crontab-go/common"
	"github.com/ywandy/crontab-go/worker"
)

var (
	confFile string //解析的配置文件路径
)

//命令行参数解析
func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "指定传入的配置文件")
	//解析
	flag.Parse()
}

//初始化一些环境
func initEnv() {
	//线程和cpu数量相同
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("初始化环境变量完成")
	fmt.Println("--> 当前程序初始化线程数量:",runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	//佛祖保佑
	common.MakeBuddhaBless()
	//程序介绍
	common.Programintro()
	//初始化命令行参数
	initArgs()
	//初始化线程
	initEnv()
	//加载配置文件
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}


ERR:
//错误的操作
	fmt.Println(err)
}
