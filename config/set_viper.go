package config

import (
	"strings"
	"fmt"
)

func (v *viper) SetSplit(key string, value []interface{}) {
	// If alias passed in, then set the proper override
	//value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := path[len(path)-1]
	fmt.Println("path:",path,"lastkey:",lastKey)
	deepestMap := deepSearch(v.config, path[0:len(path)-1])
	fmt.Println("deepestMap:",deepestMap[lastKey])

	// set innermost value
	deepestMap[lastKey] = value
}

func (v *viper) Set(key string, value interface{}) {
	// If alias passed in, then set the proper override
	//value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := path[len(path)-1]
	fmt.Println("path:",path,"lastkey:",lastKey)
	deepestMap := deepSearch(v.config, path[0:len(path)-1])
	fmt.Println("deepestMap:",deepestMap[lastKey])

	// set innermost value
	deepestMap[lastKey] = value
}

func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}