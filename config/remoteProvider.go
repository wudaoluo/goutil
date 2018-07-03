package config

import (
	"log"
	"reflect"
	"goutil/config/backends"
)

func (v *viper) AddRemoteProvider(provider, endpoint, configfile string) error {
	if !stringInSlice(provider, supportedRemoteProviders) {
		return unsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		log.Println("添加远程配置中心", provider, endpoint)
		//rp := &defaultRemoteProvider{
		//	endpoint: endpoint,
		//	provider: provider,
		//	path:     path,
		//}

		rp := &backends.ProviderConfig{
			Provider:provider,

			ConfigFiles:configfile,

		}
		//防止重复添加相同的远程配置中心
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}


func (v *viper) providerPathExists(p *backends.ProviderConfig) bool {
	for _, y := range v.remoteProviders {
		if reflect.DeepEqual(y, p) {
			return true
		}
	}
	return false
}