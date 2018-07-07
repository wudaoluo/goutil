package etcd

import (
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"fmt"
)

// Client is a wrapper around the etcd client
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

// GetValues queries etcd for keys prefixed by prefix.
//func (c *Client) GetValues(keys []string) (map[string]string, error) {
//	vars := make(map[string]string)
//	for _, key := range keys {
//		resp, err := c.client.Get(context.Background(), key, &client.GetOptions{
//			Recursive: true,
//			Sort:      true,
//			Quorum:    true,
//		})
//		if err != nil {
//			return vars, err
//		}
//		err = nodeWalk(resp.Node, vars)
//		if err != nil {
//			return vars, err
//		}
//	}
//	return vars, nil
//}

// nodeWalk recursively descends nodes, updating vars.
func nodeWalk(node *client.Node, vars map[string]string) error {
	if node != nil {
		key := node.Key
		if !node.Dir {
			vars[key] = node.Value
		} else {
			for _, node := range node.Nodes {
				nodeWalk(node, vars)
			}
		}
	}
	return nil
}


func (c *Client) watch(fn func() error) {
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

	for {
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		res, err := watcher.Next(ctx)
		if err != nil {
			//logger.Error("etcd watch 错误",err)
			fmt.Println(err)
			time.Sleep(time.Second*1)
			continue
		}

		if res.Action == "set" || res.Action == "update" || res.Action == "delete"{
			//ress := res
			fmt.Println(res.Action,res.Node.Key,res.Node.Value)
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
}


func (c *Client) WatchConfig(fn func() error) {
	go c.watch(fn)
}




func (c *Client) StopWatch() {
	if c.close {
		close(c.stopChan)
		c.close = false
	}
}

//func (c *Client) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
//	// return something > 0 to trigger a key retrieval from the store
//	if waitIndex == 0 {
//		return 1, nil
//	}
//
//	// Setting AfterIndex to 0 (default) means that the Watcher
//	// should start watching for events starting at the current
//	// index, whatever that may be.
//	watcher := c.client.Watcher(prefix, &client.WatcherOptions{AfterIndex: uint64(0), Recursive: true})
//	ctx, cancel := context.WithCancel(context.Background())
//	cancelRoutine := make(chan bool)
//	defer close(cancelRoutine)
//
//	go func() {
//		select {
//		case <-stopChan:
//			cancel()
//		case <-cancelRoutine:
//			return
//		}
//	}()
//
//	for {
//		resp, err := watcher.Next(ctx)
//		if err != nil {
//			switch e := err.(type) {
//			case *client.Error:
//				if e.Code == 401 {
//					return 0, nil
//				}
//			}
//			return waitIndex, err
//		}
//
//		// Only return if we have a key prefix we care about.
//		// This is not an exact match on the key so there is a chance
//		// we will still pickup on false positives. The net win here
//		// is reducing the scope of keys that can trigger updates.
//		for _, k := range keys {
//			if strings.HasPrefix(resp.Node.Key, k) {
//				return resp.Node.ModifiedIndex, err
//			}
//		}
//	}
//}



