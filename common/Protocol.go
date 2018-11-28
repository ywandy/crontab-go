package common

//定时任务
type Job struct {
	Name     string `json:"name"`      //任务名
	Command  string `json:"command"`   //shell命令
	CronExpr string `json:"cron_expr"` //Cron表达式
}
