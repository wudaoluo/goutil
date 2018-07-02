package config

import (
	"testing"
	"fmt"
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

	v := New()
	v.SetConfig("viper.json","json","/etc")
	v.SetDefault("debug",false)
	v.SetDefault("port1",7777)

	err := v.ReadConfig()
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
	fmt.Println(v.Getconfig())


	fmt.Println(v.Getdefault())

}
