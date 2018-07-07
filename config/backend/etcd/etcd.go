package etcd

import (
	goetcd "github.com/coreos/etcd/client"
	"time"
	"context"
	"fmt"
	"strings"
	"goutil/config/backends"
)



type Response struct {
	Action string
	Key   string
	Value []byte
	Error error
}


type client struct {
	prefix string
	keysAPI   goetcd.KeysAPI
	waitIndex uint64
}


func NewClient(endpoint []string,Prefix string) (*client,error) {
	cfg := goetcd.Config{
		Endpoints:   endpoint,
		HeaderTimeoutPerRequest: time.Second * 1,
	}

	c, err := goetcd.New(cfg)

	if err != nil {
		return nil, err
	}
	keysAPI := goetcd.NewKeysAPI(c)
	return &client{keysAPI: keysAPI,prefix:Prefix}, nil

}

/*	Get(key string) ([]byte, error)

	// List retrieves all keys and values under a provided key.
	List(key string) (KVPairs, error)

	// Set sets the provided key to value.
	Set(key string, value []byte) error

	// Watch monitors a K/V store for changes to key.
	Watch(key string, stop chan bool) <-chan *Response*/
//func (c *client) Get(key string) ([]byte,error) {
//	return nil,nil
//}
//
//
//func (c *client) Set(key, value []byte) error {
//	return nil
//}

//func (c *client) List(key string) (backend.KVPairs, error) {
//	return nil,nil
//}

func (c *client) Watch(stop chan struct{}) <-chan *backends.Response {
	respChan := make(chan *backends.Response, 10)  //加个缓冲区
	go func() {
		watcher := c.keysAPI.Watcher(c.prefix, &goetcd.WatcherOptions{
			Recursive: true,
		})
		ctx, cancel := context.WithCancel(context.TODO())

		go func() {
			<-stop
			cancel()
		}()

		respdata := &backends.Response{
			Error:nil,
		}

		for {
			var resp *goetcd.Response
			var err error
			resp, err = watcher.Next(ctx)
			if err != nil {
				respdata.Error = err
				respChan <- respdata
				time.Sleep(time.Second * 1)
				continue
			}

			respdata.Action = resp.Action

			switch resp.Action {
			case "set","update":
				respdata.Key = convertKey(resp.Node.Key,c.prefix)
				respdata.Value = resp.Node.Value

			case "delete":
				respdata.Key = convertKey(resp.Node.Key,c.prefix)

			default:
				respdata.Error = fmt.Errorf("没有发现的action")

			}

			respChan <- respdata

		}


	}()
	return respChan
}

func convertKey(key string,prefix string) string{
	a:=strings.TrimPrefix(key,prefix+"/")
	return strings.Replace(a,"/",".",-1)
}