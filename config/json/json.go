package json

import (
	"github.com/json-iterator/go"
	"io"
	"bytes"
)

type conf struct {
	configfile string
}


func NewConf() *conf {
	return &conf{}
}


func (cfg *conf) SetConfigFile(configfile string) {
	cfg.configfile = configfile
}

func (cfg *conf) ReadConfig(in io.Reader, c map[string]interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

    err := jsoniter.Unmarshal(buf.Bytes(), &c)
    return err
}


func (cfg *conf) WriteConfig() error {
	return nil
}

func (cfg *conf) WatchConfig() error {
	return nil

}
