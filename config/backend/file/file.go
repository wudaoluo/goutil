package file

import (
	"github.com/fsnotify/fsnotify"
	"goutil/config/backends"
)



type Response struct {
	Action string
	Key   string
	Value []byte
	Error error
}


type client struct {
	configFile string
}


func NewClient(configFile string) (*client,error) {
		return &client{
			configFile: configFile,
		},nil
}


func (c *client) Watch(stop chan struct{}) <-chan *backends.Response {
	respChan := make(chan *backends.Response, 10)  //加个缓冲区

	go func() {

		watcher, err := fsnotify.NewWatcher()
		//监视配置文件inode 出错了,退出程序
		if err != nil {
			panic(err)
		}

		watcher.Add(c.configFile)

		go func() {
			<-stop
			watcher.Close()
		}()

		respdata := &backends.Response{
			Error: nil,
		}

		for {
			select {
			case event := <-watcher.Events:
				//fmt.Println(event)
				if event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename ||
					event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create {
					watcher.Remove(c.configFile)
					watcher.Add(c.configFile)

					//需要读取配置文件
					//通过chan通知
					respChan <- respdata
				}

			case err := <-watcher.Errors:
				respdata.Error = err
				respChan <- respdata
			}

		}
	}()

	return respChan
}




