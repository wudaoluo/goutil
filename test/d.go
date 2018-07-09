package main

import "fmt"

func main() {
	a := make(map[string]string)
	go func() {
		for {
			a["debug"]  = "false"
		}
	}()

	for {
		fmt.Println(a["debug"])
	}
}
