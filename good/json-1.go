/*
把所有数据都放在最外层（顶层）
虽然我极其推荐你将动态可变的部分放在一个单独的 key 下面，
但是有时你可能需要处理一些预先存在的数据，它们并没有用这样的方式进行格式化。
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const input = `
{
    "type": "sound",
    "description": "dynamite",
    "authority": "the Bruce Dickinson"
}
`

type Envelope struct {
	Type string
}

type Sound struct {
	Description string
	Authority   string
}

func main() {
	var env Envelope
	buf := []byte(input)
	if err := json.Unmarshal(buf, &env); err != nil {
		log.Fatal(err)
	}
	switch env.Type {
	case "sound":
		var s struct {
			Envelope
			Sound
		}
		if err := json.Unmarshal(buf, &s); err != nil {
			log.Fatal(err)
		}
		var desc string = s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}
