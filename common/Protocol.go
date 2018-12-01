package common

import (
	"encoding/json"
	"strings"
	"github.com/gorhill/cronexpr"
	"time"
)

//常用工具类
//首字母大写才能序列化
//定时任务
type Job struct {
	Name     string `json:"name"`      //任务名
	Command  string `json:"command"`   //shell命令
	CronExpr string `json:"cron_expr"` //Cron表达式
}

//任务调度计划
//扫描下次调度时间然后判断是否执行
type JobSchedulerPlan struct {
	Job      *Job                 //属于的Job
	Expr     *cronexpr.Expression //Cron表达式
	NextTime time.Time            //下次的调度时间
}

//HTTP接口应答
type HTTPResponseContent struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

//变化事件
type JobEvent struct {
	EventType int  //save && delete事件
	Job       *Job //删除或者保存的Job内容
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

//反序列化job
func UnpackJob(value []byte) (retjob *Job, err error) {
	var (
		job *Job
	)
	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	retjob = job
	return
}

//从etcd的key中提取任务名字
//从/cron/jobs/job1 提取 job1
func ExtractJobName(jobkey string) (string) {
	return strings.Trim(jobkey, Job_Save_Dir)
}

//任务变化的事件
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

//构造执行计划
func BuildJobSchedulePlan(job *Job) (jobSchedulerPlan *JobSchedulerPlan,err  error) {
	var (
		expr *cronexpr.Expression
	)
	//解析Job的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}
	jobSchedulerPlan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()), //从当前时间算出下一次时间
	}
	return
}
