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
	"context"
	"strings"
	"encoding/json"
)

//接口添加able 组合 viperable
type viperable interface {
	//配置文件的初始化
	SetDefault(key string, value interface{})
	SetConfig(cfgfile,cfgtype string,cfgpath ...string)
	SetKeyDelim(delim string)
	ReadConfig() error
	WriteConfig() error
	WatchConfig(remoteCfg *backends.ProviderConfig) error
	//Getdefault() map[string]interface{}
	//Getconfig() map[string]interface{}

	//Get方法
	Getvalue
	//AddConfigPath(cfgpath string)
	//Operater
	//对配置文件的操作
	//ReadInConfig() error
	//WriteInConfig() error

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


	// 配置值相关的
	config         map[string]interface{}
	//override       map[string]interface{}
	defaults       map[string]interface{}
	//kvstore        map[string]interface{}
	//pflags         map[string]FlagValue
}


func New() viperable {
	v := new(viper)

	//key的分隔符，分割成列表
	v.keyDelim = "."
	v.config = make(map[string]interface{})
	v.defaults = make(map[string]interface{})

	return v
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

	return nil
}



func (v *viper) WriteConfig() error {
	return v.operating.WriteConfig()
}


//func (v *viper) WatchConfig(remoteCfg *backends.ProviderConfig) error {
//	remoteCfg.ConfigFiles = v.configFile
//	r,err := backends.New(remoteCfg)
//	if err != nil {
//		return err
//	}
//
//
//	r.WatchConfig(v)
//
//	return nil
//}

func (v *viper) WatchConfig(remoteCfg *backends.ProviderConfig) error {
	c ,_:=NewClient(remoteCfg.Endpoint,remoteCfg.Prefix)
	watcher := c.client.Watcher(c.prefix, &client.WatcherOptions{
		Recursive: true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	//确保整个for完全退出后在 退出这个函数，说白了就是善后工作
	cancelRoutine := make(chan bool)
	defer close(cancelRoutine)

	go func() {
		select {
		case <-c.stopChan:
			cancel()
		case <-cancelRoutine:
			return
		}
	}()

	go func() {

		for {
			res, err := watcher.Next(ctx)
			if err != nil {
				//logger.Error("etcd watch 错误",err)
				fmt.Println(err)
				time.Sleep(time.Second * 1)
				continue
			}

			if res.Action == "set" || res.Action == "update" || res.Action == "delete" {
				//ress := res

				//aa := []string{}
				//jsonStringToObject(res.Node.Value,&aa)
				if strings.HasPrefix(res.Node.Value,"[") &&
					strings.HasSuffix(res.Node.Value,"]") {
					a := []string{}
					data := []byte(res.Node.Value)
					json.Unmarshal(data, a)
					v.config["key3"] = a
					return
				}


				v.config["key3"] = res.Node.Value
				return
				//v.config["key4"] = strings.Split(res.Node.Value,",")
				fmt.Println(res.Action, res.Node.Key, res.Node.Value)
				//go func() {
				//logger.Info("更新",res.Node.Key,res.Node.Value)

				//_dir,_key,err := getkey(ress.Node.Key)
				//if err != nil {
				//logger.Error(err)
				//}
				//更新值 ，保存到配置文件中

				//cfgfile.save()
				//}()
			}

		}
	}()
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

//func (v *viper) WatchConfig() error {
//	fn := func() {
//		watcher, err := fsnotify.NewWatcher()
//		if err != nil {
//			log.Fatal(err)
//		}
//		defer watcher.Close()
//
//		watcher.Add(v.configFile)
//
//		for {
//			select {
//			case event := <-watcher.Events:
//				fmt.Println("event:", event)
//
//				if event.Op&fsnotify.Remove == fsnotify.Remove ||
//					event.Op&fsnotify.Rename == fsnotify.Rename ||
//					event.Op&fsnotify.Write == fsnotify.Write ||
//					event.Op&fsnotify.Create == fsnotify.Create {
//					watcher.Remove(v.configFile)
//					watcher.Add(v.configFile)
//					err := v.ReadConfig()
//					if err != nil {
//						log.Println("error:", err)
//					}
//				}
//
//			case err := <-watcher.Errors:
//				log.Println("error:", err)
//			}
//		}
//
//	}
//
//	go fn()
//
//	return nil
//}

func (v *viper) Getdefault() map[string]interface{} {
	return v.defaults
}

func (v *viper) Getconfig() map[string]interface{} {
	return v.config
}

