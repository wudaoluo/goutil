package config

import "fmt"


//表示遇到不受支持配置格式
type UnsupportedConfigError string


func (e UnsupportedConfigError) Error() string {
	return fmt.Sprintf("这是一个不支持的配置类型 %q", string(e))
}


type NotFoundConfigError string

func (e NotFoundConfigError) Error() string {
	return fmt.Sprintf("没有找到配置文件 %s", string(e))
}


type configReadError string

func (e configReadError) Error() string {
	return fmt.Sprintf("配置文件读取失败 %s", string(e))
}


// ConfigParseError denotes failing to parse configuration file.
type configParseError struct {
	err error
}

// Error returns the formatted configuration error.
func (pe configParseError) Error() string {
	return fmt.Sprintf("While parsing config: %s", pe.err.Error())
}