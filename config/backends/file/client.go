package file

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
)



//var replacer = strings.NewReplacer("/", "_")

// Client provides a shell for the yaml client
type Client struct {
	configFile string
	filter     string
	stopChan   chan struct{}
	close      bool
}

type ResultError struct {
	response uint64
	err      error
}

func NewClient(configFile string, filter string) (*Client,error) {
		return &Client{
			configFile: configFile,
			filter: filter,
			stopChan: make(chan struct{}),
			close:true},nil
}


func (c *Client) watch(fn func() error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	watcher.Add(c.configFile)

	for {
		select {
		case event := <-watcher.Events:
			fmt.Println("event:", event)

			if event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				watcher.Remove(c.configFile)
				watcher.Add(c.configFile)

				//需要读取配置文件
				err := fn()
				if err != nil {
					log.Println("error:", err)
				}
			}

		case err := <-watcher.Errors:
			log.Println("error:", err)

		case <- c.stopChan:
			log.Println("关闭config watch")
			watcher.Close()
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


