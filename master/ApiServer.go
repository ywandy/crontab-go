package master

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
	"strconv"
	"github.com/ywandy/crontab-go/common"
	"encoding/json"
)

type ApiServer struct {
	httpServer *http.Server
}

//单例对象的全局访问标识
var (
	G_ApiServer *ApiServer
)

//提交任务的接口

func handleJobSave(ctx *gin.Context) {
	var (
		postjob  string //接收post表单的job字段
		respbyte []byte
		err      error
		job      common.Job
		oldjob   *common.Job
	)
	//解析表单
	postjob = ctx.PostForm("job")
	//判断表单
	if postjob == "" {
		respbyte = common.MakeResponse(-1, "Post错误", "参数不能为空")
		ctx.String(200, string(respbyte))
	}
	//处理表单
	//反序列化job
	if err = json.Unmarshal([]byte(postjob), &job); err != nil {
		goto ERR
	}
	if oldjob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//返回值
	respbyte = common.MakeResponse(0, "成功", oldjob)
	ctx.String(200, string(respbyte))
	return
ERR:
//致命错误弹出
	respbyte = common.MakeResponse(-1, "致命错误", string(err.Error()))
	ctx.String(200, string(respbyte))
	return
}

//删除任务的接口
// post {name=job1}
func handleJobDelete(ctx *gin.Context) {
	var (
		postjob  string //接收post表单的job字段
		respbyte []byte
		err      error
		job      common.Job
		oldjob   *common.Job
	)
	//解析表单
	postjob = ctx.PostForm("name")
	//判断表单
	if postjob == "" {
		respbyte = common.MakeResponse(-1, "Post错误", "参数不能为空")
		ctx.String(200, string(respbyte))
	}
	//处理表单
	//反序列化job
	if err = json.Unmarshal([]byte(postjob), &job); err != nil {
		goto ERR
	}
	if oldjob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//返回值
	respbyte = common.MakeResponse(0, "成功", oldjob)
	ctx.String(200, string(respbyte))
	return
ERR:
//致命错误弹出
	respbyte = common.MakeResponse(-1, "致命错误", string(err.Error()))
	ctx.String(200, string(respbyte))
	return
}


func InitApiServer() (err error) {
	var (
		router *gin.Engine
	)
	//release 模式
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	//注册路由
	router.POST("/job/save", handleJobSave)
	router.POST("/job/delete",handleJobDelete)
	//使用协程去启动
	go router.Run(":" + strconv.Itoa(G_Config.API_PORT))
	fmt.Println("web 服务器已经运行在", ":"+strconv.Itoa(G_Config.API_PORT))
	return
}
