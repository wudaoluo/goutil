package config

import (
	"testing"
	"fmt"
	"time"
)

/*

func Test_JsonConf(t *testing.T) {
	a := GetInstance()
	a.ParseConf("../cmd/server.json","json")
	fmt.Println(a.Cfg.Version)
	a.Cfg.Version = "1.121"
	a.Cfg.Mysql.DBpasswd = "aa"
	fmt.Println(a.Cfg.Version)

	a.C.SaveFile()
	fmt.Println(a.Cfg.Version)
	time.Sleep(10*time.Second)
	a.C.Reload()
	fmt.Println(a.Cfg.Version)
	t.Log(a.Cfg)
}

*/

func Test_viper(t *testing.T) {
	var err error
	v := New()
	v.SetConfig("viper.json","json","/etc")
	v.SetDefault("debug",false)
	v.SetDefault("port1",7777)

	//开启动态配置文件
	err = v.WatchConfig()
	//err = v.AddRemoteProvider("etcd", "http://127.0.0.1:4001","/config/hugo.json")
	if err != nil {
		panic(err)
	}
	err = v.ReadConfig()
	//err := v.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}


	fmt.Println(v.GetBool("debug"))
	fmt.Println(v.GetString("port"))
	fmt.Println(v.GetInt("port1"))
	fmt.Println(v.GetString("port1"))
	fmt.Println(v.GetString("a.test"))
	fmt.Println(v.GetString("a.c.as"))
	fmt.Println(v.GetString("a.c.cb.ca"))
	fmt.Println(v.GetStringMap("a.c.cb")["ca"])

	for {
		time.Sleep(1*time.Second)
		fmt.Println(v.GetBool("debug"))
	}
	//fmt.Println(v.Getconfig())


	//fmt.Println(v.Getdefault())



}
