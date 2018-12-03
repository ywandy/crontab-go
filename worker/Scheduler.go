package worker

import (
	"github.com/ywandy/crontab-go/common"
	"time"
	"fmt"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent               //etcd任务事件队列，使用chan实现的
	jobPlanTable map[string]*common.JobSchedulerPlan //内存的任务计划表
}

var (
	G_scheduler *Scheduler
)

//处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulerPlan
		err             error
		jobExisted      bool
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE: //删除任务事件
		//因为ETCD即使没有的条目，发送删除事件，依旧会推送过来，因此先要判断内存的map中是否存在这样的一个条目
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}


func (scheduler *Scheduler)TryStartJob(jobPlan *common.JobSchedulerPlan){

}

//重新计算任务调度的状态
//返回下次调度时间间隔
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulerPlan
		now      time.Time
		nearTime *time.Time
	)
	//如果任务表为空，责睡眠多久
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}
	//获取当前的时间
	now = time.Now()
	//1. 遍历所有任务
	//2.过期的任务立即执行
	//3.统计最近要过期的任务时间（N秒后过去 == scheduleAfter
	//下次调度的间隔（最近的任务调度时间-当前时间）
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO 尝试执行任务
			//因为前一次还没结束，这次不启动

			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
			fmt.Println("任务名字:",jobPlan.Job.Name,"当前时间:",now,"下次执行时间",jobPlan.NextTime)
		}
		//统计最近要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	//返回调度间隔时间
	scheduleAfter = (*nearTime).Sub(now)
	return
}

//调度循环协程
func (scheduler *Scheduler) shcedulerLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)
	//初始化一次
	scheduleAfter = scheduler.TrySchedule()
	//调度的延迟定时器
	scheduleTimer = time.NewTimer(scheduleAfter)
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			//对内存维护的任务列表做增删改查
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: //任务到期
		}
		//调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		//重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

//推送任务变化的事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent //把接收的jobEvent发送到队列
}

//不断循环
//初始化调度器
func InitScheduler() (err error){
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000), ///队列能存1000个数值
		jobPlanTable:make(map[string]*common.JobSchedulerPlan),
	}
	go G_scheduler.shcedulerLoop() //使用协程去启动事件的主循环
	return
}
