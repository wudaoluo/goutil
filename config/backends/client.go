package backends

import (
	"errors"
)


// RemoteProvider 用来连接到远程配置中心，初始化
//type RemoteProvider interface {
//	Provider() string
//	Endpoint() string
//	Path() string
//}

type StoreClient interface {
	GetValues(keys []string) (map[string]string, error)
	WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

// New is used to create a storage client based on our configuration.
func New(config ProviderConfig) (StoreClient, error) {
	switch config.Provider {
	case "etcd":
		// Create the etcd client upfront and use it for the life of the process.
		// The etcdClient is an http.Client and designed to be reused.
		//return etcd.NewEtcdClient(endpoint)

	case "file":
		//return file.NewFileClient(config.ConfigFiles,config.Filter)

	case "consul":
	case "etcdv3":
	case "redis":
		
	}
	return nil, errors.New("Invalid backend")
}
