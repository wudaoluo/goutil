package config

import (
	"strings"
	"goutil/config/json"
	"goutil/config/ini"
	"io"
)


type Operater interface {
	SetConfigFile(cnofigfile string)  //设置配置文件
	ReadConfig(in io.Reader, c map[string]interface{}) error   //读取配置文件
	WriteConfig(v interface{}) error  //写入配置文件
	//WatchConfig() error
	//WatchRemoteConfig()
}


//type operating struct {}

//查找配置文件类型 并初始化
func getType(cfgtype string) (Operater,error) {
	//转换成小写
	cfgtype = strings.ToLower(cfgtype)

	if !stringInSlice(cfgtype,supportedExts) {
		return nil,UnsupportedConfigError(cfgtype)
	}

	var a Operater

	switch cfgtype {
	case "json":
		a = json.NewConf()
	case "ini":
		a = ini.NewConf()
	}

	return a,nil
}




/*
	switch strings.ToLower(v.getConfigType()) {
	case "yaml", "yml":
		if err := yaml.Unmarshal(buf.Bytes(), &c); err != nil {
			return ConfigParseError{err}
		}

	case "json":
		if err := json.Unmarshal(buf.Bytes(), &c); err != nil {
			return ConfigParseError{err}
		}

	case "hcl":
		obj, err := hcl.Parse(string(buf.Bytes()))
		if err != nil {
			return ConfigParseError{err}
		}
		if err = hcl.DecodeObject(&c, obj); err != nil {
			return ConfigParseError{err}
		}

	case "toml":
		tree, err := toml.LoadReader(buf)
		if err != nil {
			return ConfigParseError{err}
		}
		tmap := tree.ToMap()
		for k, v := range tmap {
			c[k] = v
		}

	case "properties", "props", "prop":
		v.properties = properties.NewProperties()
		var err error
		if v.properties, err = properties.Load(buf.Bytes(), properties.UTF8); err != nil {
			return ConfigParseError{err}
		}
		for _, key := range v.properties.Keys() {
			value, _ := v.properties.Get(key)
			// recursively build nested maps
			path := strings.Split(key, ".")
			lastKey := strings.ToLower(path[len(path)-1])
			deepestMap := deepSearch(c, path[0:len(path)-1])
			// set innermost value
			deepestMap[lastKey] = value
		}
	}
*/
//
//func (o *operating) ReadConfig() error {
//	return nil
//}
//
//
//func (o *operating) WriteConfig() error {
//	return nil
//}