package master

import (
	"net/http"
	"net"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

//单例对象的全局访问标识
var (
	G_ApiServer *ApiServer
)

//保存任务的回调函数
func handleJobSave(w http.ResponseWriter, r *http.Request) {

}

//初始化http服务(大写能被其他模块调用)
func InitApiServer() (err error){
	//配置路由
	var (
		mux      *http.ServeMux
		httpListener net.Listener
		httpServer *http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	if httpListener, err = net.Listen("tcp", ":8070"); err != nil {
		return
	}
	//创建http server
	httpServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
	}
	G_ApiServer = &ApiServer{
		httpServer:httpServer,
	}
	go httpServer.Serve(httpListener)
	return nil
}
