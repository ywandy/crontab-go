package master

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

//配置文件的结构体
type Config struct {
	API_PORT          int      `json:"api_port"`
	API_Read_Timeout  int      `json:"api_read_timeout"`
	API_Write_Timeout int      `json:"api_write_timeout"`
	EtcdEndpoints     []string `json:"etcdEndpoints"`
	EtcdDialTimeout   int      `json:"etcdDialTimeout"`
}

var (
	G_Config *Config
)

//配置文件的init方法
//传入配置文件的路径
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	//读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	//json反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	//赋值单例
	G_Config = &conf
	fmt.Println("载入配置文件完成")
	return
}
