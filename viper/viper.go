package viper

import (
	"strings"
	"github.com/magiconair/properties"
	"github.com/fsnotify/fsnotify"
)

type Viper struct {
	//分隔键列表的分隔符
	keyDelim string

	// A set of paths to look for the config file in
	configPaths []string



	// A set of remote providers to search for the configuration
	remoteProviders []*defaultRemoteProvider

	// Name of file to look for inside the path
	configName string
	configFile string
	configType string
	envPrefix  string

	automaticEnvApplied bool
	envKeyReplacer      *strings.Replacer

	config         map[string]interface{}
	override       map[string]interface{}
	defaults       map[string]interface{}
	kvstore        map[string]interface{}
	pflags         map[string]FlagValue
	env            map[string]string
	aliases        map[string]string
	typeByDefValue bool

	// Store read properties on the object so that we can write back in order with comments.
	// This will only be used if the configuration read is a properties file.
	properties *properties.Properties

	onConfigChange func(fsnotify.Event)
}



func New() *Viper {
	v := new(Viper)
	v.keyDelim = "."
	v.configName = "config"
	v.config = make(map[string]interface{})
	v.override = make(map[string]interface{})
	v.defaults = make(map[string]interface{})
	v.kvstore = make(map[string]interface{})
	v.pflags = make(map[string]FlagValue)
	v.env = make(map[string]string)
	v.aliases = make(map[string]string)
	v.typeByDefValue = false

	return v
}


