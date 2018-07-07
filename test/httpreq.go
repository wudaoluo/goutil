package main

import (
	"time"
	"net/http"
	"net"
	"io/ioutil"
	"fmt"
)


func main() {
	c := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Duration(3) * time.Second)
				c, err := net.Dial(netw, "114.55.73.227:80")
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	resp,err := c.Get("http://down.163.com/script/init/install.sh")

	/*
	大多数情况下，当你的http响应失败时，resp变量将为nil，而err变量将是non-nil。
	然而，当你得到一个重定向的错误时，两个变量都将是non-nil。这意味着你最后依然会内存泄露。
	所以要这样写
	*/

	if resp != nil {
		defer resp.Body.Close()
	}

	if err!= nil {
		fmt.Println(err)
		return
	}

	r,err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(r))
}
