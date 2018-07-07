package backends

import (
	"errors"
	"goutil/config/backends/file"
	"goutil/config/backends/etcd"

)


// RemoteProvider 用来连接到远程配置中心，初始化
//type RemoteProvider interface {
//	Provider() string
//	Endpoint() string
//	Path() string
//}

type StoreClient interface {
	WatchConfig(func() error)
	StopWatch()
	//GetValues(keys []string) (map[string]string, error)
	//WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

// New is used to create a storage client based on our configuration.
func New(config *ProviderConfig) (StoreClient,error) {
	if config.Provider == "" {
		config.Provider = "file"
	}
	switch config.Provider {
	case "etcd":
		// Create the etcd client upfront and use it for the life of the process.
		// The etcdClient is an http.Client and designed to be reused.
		return etcd.NewClient(config.Endpoint,config.Prefix)

	case "file":
		return file.NewClient(config.ConfigFiles,config.Filter)

	case "consul":
	case "etcdv3":
	case "redis":
		
	}
	return nil, errors.New("无效的backend")
}
