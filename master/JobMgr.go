package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"github.com/ywandy/crontab-go/common"
	"encoding/json"
	"context"
	"fmt"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

//单例
var (
	G_jobMgr *JobMgr
)

//初始化job管理器
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		lease  clientv3.Lease
		kv     clientv3.KV
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
	//赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	fmt.Println("初始化jobMgr完成")
	return
}

//实现jobmgr保存
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	//把任务保存 /cron/jobs/任务名 -> json
	var (
		jobKey    string
		jobVal    []byte
		putRespon *clientv3.PutResponse
		oldJobObj common.Job
	)
	//保存key
	jobKey = common.Job_Save_Dir + job.Name
	//序列化json得到jobval
	if jobVal, err = json.Marshal(job); err != nil {
		return
	}
	//保存etcd
	//保存成功返回旧值
	if putRespon, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobVal), clientv3.WithPrevKV()); err != nil {
		return
	}
	//如果是更新，那么是返回旧的值
	if putRespon.PrevKv != nil {
		//对旧值做反序列化
		if err = json.Unmarshal(putRespon.PrevKv.Value, &oldJobObj); err != nil {
			err = nil //不需要得到旧值是否合法
			return
		}
		//返回旧值
		oldJob = &oldJobObj
	}
	return
}

//实现jobmgr删除
//返回删除结果
func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	//把任务保存 /cron/jobs/任务名 -> json
	var (
		jobKey         string
		deleteResponse *clientv3.DeleteResponse
		oldJobObj      common.Job
	)
	//得到要删除的任务的key
	jobKey = common.Job_Save_Dir + name
	//删除key
	if deleteResponse, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	//返回被删除的信息
	if len(deleteResponse.PrevKvs) != 0 {
		//解析旧值
		if err = json.Unmarshal(deleteResponse.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil //无论旧值是否能解析，只要能删除就好
			return
		}
	}
	oldJob = &oldJobObj
	return
}
