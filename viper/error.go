package viper

import "fmt"

// UnsupportedConfigError denotes encountering an unsupported
// configuration filetype.
type UnsupportedConfigError string

// Error returns the formatted configuration error.
func (str UnsupportedConfigError) Error() string {
	return fmt.Sprintf("Unsupported Config Type %q", string(str))
}

// UnsupportedRemoteProviderError denotes encountering an unsupported remote
// provider. Currently only etcd and Consul are supported.
type UnsupportedRemoteProviderError string

// Error returns the formatted remote provider error.
func (str UnsupportedRemoteProviderError) Error() string {
	return fmt.Sprintf("Unsupported Remote Provider Type %q", string(str))
}

// RemoteConfigError denotes encountering an error while trying to
// pull the configuration from the remote provider.
type RemoteConfigError string

// Error returns the formatted remote provider error
func (rce RemoteConfigError) Error() string {
	return fmt.Sprintf("Remote Configurations Error: %s", string(rce))
}

// ConfigFileNotFoundError denotes failing to find configuration file.
type ConfigFileNotFoundError struct {
	name, locations string
}

// Error returns the formatted configuration error.
func (fnfe ConfigFileNotFoundError) Error() string {
	return fmt.Sprintf("Config File %q Not Found in %q", fnfe.name, fnfe.locations)
}
