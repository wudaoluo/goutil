package json

import (
	"github.com/json-iterator/go"
	"io"
	"bytes"
	"io/ioutil"
	"os"
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


func (cfg *conf) WriteConfig(v interface{}) error {
	var err error

	saveData,err := jsoniter.MarshalIndent(v,"","    ")
	if err != nil {
		return err
	}

	temp, err := ioutil.TempFile(".", ".tmp")
	if err != nil {
		return err
	}
	defer temp.Close()

	_, err = temp.Write(saveData)
	if err != nil {
		return err
	}

	defer os.Remove(temp.Name())


	err = os.Rename(temp.Name(),cfg.configfile)

	if err != nil {
		return err
	}

	return nil
}
