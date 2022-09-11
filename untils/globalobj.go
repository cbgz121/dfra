package untils

import (
	"dfra/diface"
	"encoding/json"
	"io/ioutil"
	"time"
)

//存储一切关于dfra的全局参数，供其他模块使用
//一些参数是可以通过dfra.json由用户进行配置

const (
	DefaultIPVersion          = "tcp4"
	DefaultHost               = "127.0.0.1"
	DefaultTcpPort            = 9999
	DefaultMaxConnSize        = 100
	DefaultMaxPackageSize     = 1024 * 1024
	DefaultWorkerPoolSize     = 10
	DefaultHeartRateInSecond  = 30 * time.Second
	DefaultHeartFreshLevel    = 5
	DefaultHeartBeatPackageId = 100
)

type Config struct {
	//Server
	TcpServer      diface.IServer //当前dfra全局的Server对象
	Host           string         //当前服务器主机监听的IP
	TcpPort        int            //当前服务器主机监听的端口号
	Name           string         //当前服务器的名称
	Version        string         //当前dfra的版本号
	MaxConn        int            //当前服务器主机允许的最大连接数
	MaxPackageSize uint32         //当前dfra框架数据包的最大值

	WorkerPoolSize     uint32 //当前业务工作Worker池的Goroutine数量
	HeartRateInSecond  time.Duration
	HeartFreshLevel    uint32
	HeartBeatPackageId uint32
}

// 定义一个全局的对外Gloablobj
var GlobalObject *Config

// 从dfra.json去加载用于自定义的参数
func (g *Config) Reload() {
	data, err := ioutil.ReadFile("/root/goproject/dfra/test/dfrademo/conf/dfra.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {

	//如果全局配置没有加载，默认的值
	GlobalObject = &Config{
		Name:               "dfraServerApp",
		Version:            "V0.4",
		TcpPort:            DefaultTcpPort,
		Host:               DefaultHost,
		MaxConn:            DefaultMaxConnSize,
		MaxPackageSize:     DefaultMaxPackageSize,
		WorkerPoolSize:     DefaultWorkerPoolSize,
		HeartRateInSecond:  DefaultHeartRateInSecond,
		HeartFreshLevel:    DefaultHeartFreshLevel,
		HeartBeatPackageId: DefaultHeartBeatPackageId,
	}

	GlobalObject.Reload()
}
