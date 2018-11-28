package common

import (
	"encoding/json"
)

//常用工具类
//首字母大写才能序列化
//定时任务
type Job struct {
	Name     string `json:"name"`      //任务名
	Command  string `json:"command"`   //shell命令
	CronExpr string `json:"cron_expr"` //Cron表达式
}

//HTTP接口应答
type HTTPResponseContent struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

//制作response的工具类
func MakeResponse(erron int, msg string, data interface{}) (respon []byte) {
	//定义一个response
	var (
		httpresponsecontent HTTPResponseContent
	)
	httpresponsecontent.Errno = erron
	httpresponsecontent.Msg = msg
	httpresponsecontent.Data = data
	//序列化json
	respon, _ = json.Marshal(httpresponsecontent)
	return
}
