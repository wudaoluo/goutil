package backend

import "goutil/config/backends"

type Config struct {
	Backend 		string
	Prefix   		string
	Endpoint 		[]string
	ConfigFiles 	string
	Fn      		func(config *backends.Response) error
	//Filter       string
}
