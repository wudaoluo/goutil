package elastic

import (
	"testing"
	"fmt"
)


func Test_JsonConf(t *testing.T) {
	es:=New("http://114.55.73.227:9200","test")
	err := es.Init()
	if err != nil {
		fmt.Println("init",err)
		return
	}

	err = es.Index()
	if err != nil {
		fmt.Println("index",err)
		return
	}

	err = es.Put("aaa",[]byte(`{"aaaa":"asas","time":"2018.06.28 15:33:33"}`))
	if err != nil {
		fmt.Println("put",err)
		return
	}
}
