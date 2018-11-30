package worker

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"context"
	"github.com/ywandy/crontab-go/common"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

//单例
var (
	G_jobMgr *JobMgr
)

//初始化job管理器
func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		lease   clientv3.Lease
		kv      clientv3.KV
		watcher clientv3.Watcher
	)
	//初始化etcd配置
	config = clientv3.Config{
		Endpoints:   G_Config.EtcdEndpoints, //集群地址
		DialTimeout: time.Duration(G_Config.EtcdDialTimeout) * time.Millisecond,
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	//得到KV和lease的api子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)
	//赋值单例
	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}
	fmt.Println("正在启动任务监听")
	G_jobMgr.watchJobs()
	fmt.Println("初始化jobMgr完成")
	return
}

//监听任务的变化
func (JobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)
	//1.get 下所有任务，并且获取当前集群的revision
	if getResp, err = JobMgr.kv.Get(context.TODO(), common.Job_Save_Dir, clientv3.WithPrefix()); err != nil {
		return
	}
	//当前有哪些任务
	for _, kvpair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			//TODO推送给调度器
		}
	}
	//2.从revision向后监听事件变化
	//监听协程
	go func() {
		//从GET后下一个版本开始监听
		watchStartRevision = getResp.Header.Revision + 1
		watchChan = JobMgr.watcher.Watch(context.TODO(), common.Job_Save_Dir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: //任务保存事件
					//反序列化etcd的key
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						err = nil
						continue
					} //非法的任务,忽略
					//构建一个保存的Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
					fmt.Println(*jobEvent)
					// TODO: 推给scheduler

				case mvccpb.DELETE: //任务被删除
					//因为我们拿到的是一个key，那么我们需要把job的名字提取出来
					//因此需要一个公用函数
					//提取任务名
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					//构建一个删除的Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
					//TODO:推送scheduler一个删除事件
					G_scheduler.PushJobEvent(jobEvent)
					fmt.Println(*jobEvent)
				}
			}
		}
	}()
	return
}
