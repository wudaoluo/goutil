package config

import (
	"log"
	"fmt"
	"path"
	"os"
	"goutil/files"
	"bytes"
	"goutil/config/backends"
	"github.com/coreos/etcd/client"
	"time"
	"goutil/config/backend"
	"sync"
)

//接口添加able 组合 viperable
type viperable interface {
	//配置文件的初始化
	SetDefault(key string, value interface{})
	SetConfig(cfgfile,cfgtype string,cfgpath ...string)
	SetKeyDelim(delim string)

	WatchConfig(remoteCfg *backend.Config) error
	Stop()
	//Getdefault() map[string]interface{}
	//Getconfig() map[string]interface{}


	//Get方法
	Getvalue

	//AddConfigPath(cfgpath string)
	//Operater

	//对配置文件的操作
	ReadConfig() error
	WriteConfig() error

	ReadRemoteConfig() error

	SetFunc(fn func(key string, value interface{}) error)

	//远程配置
	//AddRemoteProvider(provider, endpoint, path string) error

}

type viper struct {
	// Delimiter that separates a list of keys
	// used to access a nested value in one go
	keyDelim string

	// A set of paths to look for the config file in
	configPaths []string

	configFile string
	configType string

	remoteProviders []*backends.ProviderConfig
	operating Operater
	backend    backend.StoreClient

	// 配置值相关的
	config         map[string]interface{}
	//override       map[string]interface{}
	defaults       map[string]interface{}
	//kvstore        map[string]interface{}
	//pflags         map[string]FlagValue
	stop       chan struct{}
	//TODO改怎么解释这个函数的作用大呢
	fn      		func(key string, value interface{}) error

	mu       sync.RWMutex
}


func New() viperable {
	v := new(viper)

	//key的分隔符，分割成列表
	v.keyDelim = "."
	v.config = make(map[string]interface{})
	v.defaults = make(map[string]interface{})
	v.stop = make(chan struct{},0)

	return v
}


func (v *viper) SetFunc(fn func(key string, value interface{}) error) {
	v.fn = fn
}

func (v *viper) Stop() {
	v.stop <- struct{}{}
}

// SetConfigName 设置配置文件的名称
func (v *viper) SetConfig(cfgfile ,cfgtype string, cfgpath ...string) {
	log.Println(fmt.Sprintf("配置文件:%s",cfgfile))

	if cfgtype == "" {
		panic("配置类型不能为空")
	}

	if cfgfile == ""  {
		panic("配置文件不能为空")
	}

	var configfile string

	var err error
	// 循环path 判断配置文件 存在(IsExist true) 赋值 configfile = _cfgfile,
	// 找不到则configfile:=cfgfile
	cfgpath = append(cfgpath,".")
	for _,i := range cfgpath {
		_cfgfile := path.Join(i,cfgfile)
		_, err := os.Stat(_cfgfile)
		if err == nil {   //TODO 使用os.IsExist(err) 每次都是都返回 false
			configfile = _cfgfile
			continue

		}
	}

	//根据类型初始化配置
	v.operating,err =getType(cfgtype)
	if err != nil {
		panic(err)
	}

	if configfile == "" {
		panic(NotFoundConfigError(cfgfile))
	}
	v.configType = cfgtype
	v.configFile = configfile


	//赋值给Operater接口
	v.operating.SetConfigFile(configfile)

}


func (v *viper) SetKeyDelim(delim string) {
	log.Println(fmt.Sprintf("设置key的分隔符 %s",delim))
	if delim == "" {
		v.keyDelim = delim
	}
}

// SetDefault 注册默认值
func (v *viper) SetDefault(key string,value interface{}) {
	v.defaults[key] = value
	log.Println(fmt.Sprintf("key:%s,value:%s",key,value))
}


//

//TODO 这里读取配置文件失败难道不应该 panic吗
func (v *viper) ReadConfig() error {
	file, err :=files.ReadFile(v.configFile)
	if err != nil {
		return configReadError(v.configFile)
	}

	err = v.operating.ReadConfig(bytes.NewReader(file),v.config)
	if err != nil {
		return configParseError{err}
	}

	for {
		v.Set("debug",false)
	}


	return nil
}



func (v *viper) ReadRemoteConfig() error {
	//读取远程配置
	var err error
	respChan := make(chan *backends.Response, 10)  //添加缓存 让etcd尽快处理完

	go func() {
		for {
			select {
			case a,ok := <- respChan:
				if !ok {
					//读取完成关闭通道，写入配置文件到本地
					err = v.WriteConfig()
					if err != nil {
						fmt.Println(err,"关闭通道")
					}
					return
				}
				err = v.Set(a.Key,a.Value)
				//if err != nil {
				//	fmt.Println(err)
				//}
				//
				//fmt.Println(v.config)

			}

		}
	}()

	err = v.backend.List(respChan)
	if err != nil {

		fmt.Println(err)
		return err
	}


	return nil
	//写入到本地
	//return v.WriteConfig()
}


func (v *viper) WriteConfig() error {
	return v.operating.WriteConfig(v.config)
}

//这里的stop 以后可以换成context,,  string, map，[]string  []int 都能自动转换,其他类型需要自定义fn函数
func (v *viper) WatchConfig(remoteCfg *backend.Config) error {
	remoteCfg.ConfigFiles = v.configFile
	var err error
	v.backend,err = backend.New(remoteCfg)
	if err != nil {
		fmt.Println(err)
		return err
	}


	remotechan := v.backend.Watch(v.stop)
	//var a backends.Response
	go func() {
		for {
			select {
			case a := <- remotechan:  //TODO 还有DELETE类型没有判断
				if a.Error != nil {
					fmt.Println("err1",err)
					continue
				}

				if remoteCfg.Backend == "file" {
					v.ReadConfig()
					continue
				}


				//这里是其他的backend
				fmt.Println(a.Key,string(a.Value.(string)))


				v.Set(a.Key,a.Value)

				//保存到本地
				err = v.WriteConfig()
				if err != nil {
					fmt.Println(err)
				}


			case <- v.stop:
				return
			}
		}
	}()
	//r.WatchConfig(v)

	return nil
}





type Client struct {
	client client.KeysAPI
	prefix string
	stopChan   chan struct{}
	close      bool
}

// NewEtcdClient returns an *etcd.Client with a connection to named machines.
func NewClient(endpoint []string,Prefix string) (*Client, error) {
	cfg := client.Config{
		Endpoints:   endpoint,
		Transport:  client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 1,
	}


	var kapi client.KeysAPI

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	kapi = client.NewKeysAPI(c)
	return &Client{client:kapi,prefix:Prefix}, nil
}



func (v *viper) Getdefault() map[string]interface{} {
	return v.defaults
}

func (v *viper) Getconfig() map[string]interface{} {
	return v.config
}

