package backend

import (
	"errors"
	"goutil/config/backend/etcd"
	"goutil/config/backend/file"
	"goutil/config/backends"
)

//type KVPair struct {
//	key string
//	value []byte
//}
//
//
//type KVPairs []*KVPair


// Response represents a response from a backend store.


// A Store is a K/V store backend that retrieves and sets, and monitors
// data in a K/V store.
type StoreClient interface {
	// Get retrieves a value from a K/V store for the provided key.
	//Get(key string) ([]byte, error)
	//
	//// List retrieves all keys and values under a provided key.
	List(respChan chan *backends.Response) error
	//
	//// Set sets the provided key to value.
	//Set(key string, value []byte) error

	// Watch monitors a K/V store for changes to key.
	Watch(stop chan struct{}) <-chan *backends.Response
}


//type Response struct {
//	Action string
//	Key   string
//	Value []byte
//	Error error
//}



func New(config *Config) (StoreClient,error) {
	if config.Backend == "" {
		config.Backend = "file"
	}
	switch config.Backend {
	case "etcd":
		// Create the etcd client upfront and use it for the life of the process.
		// The etcdClient is an http.Client and designed to be reused.
		//return etcd.NewClient(config.Endpoint,config.Prefix)
		return etcd.NewClient(config.Endpoint,config.Prefix)

	case "file":
		return file.NewClient(config.ConfigFiles)

	case "consul":
	case "etcdv3":
	case "redis":

	}
	return nil, errors.New("无效的backend")
}

