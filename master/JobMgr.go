package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
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
	return
}
