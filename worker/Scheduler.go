package worker

import "github.com/ywandy/crontab-go/common"

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
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name];jobExisted{
			delete(scheduler.jobPlanTable,jobEvent.Job.Name)
		}
	}
}

//调度循环协程
func (scheduler *Scheduler) shcedulerLoop() {
	var (
		jobEvent *common.JobEvent
	)
	for {
		select {
		case <-scheduler.jobEventChan:
			//对内存维护的任务列表做增删改查
		}
	}
}

//推送任务变化的事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent //把接收的jobEvent发送到队列
}

//不断循环
//初始化调度器
func InitScheduler() {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000), ///队列能存1000个数值
	}
	go G_scheduler.shcedulerLoop() //使用协程去启动事件的主循环
	return
}
