package config

import "io"

type defaultRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
}

func (rp defaultRemoteProvider) Provider() string {
	return rp.provider
}

func (rp defaultRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp defaultRemoteProvider) Path() string {
	return rp.path
}


// RemoteProvider 用来连接到远程配置中心，初始化
type RemoteProvider interface {
	Provider() string
	Endpoint() string
	Path() string
}


//remoteConfigFactory 对远程配置中心的操作
type remoteConfigFactory interface {
	Get(rp RemoteProvider) (io.Reader, error)
	Watch(rp RemoteProvider) (io.Reader, error)
	WatchChannel(rp RemoteProvider) (<-chan *RemoteResponse, chan bool)
}

type RemoteResponse struct {
	Value []byte
	Error error
}