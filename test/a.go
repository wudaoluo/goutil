package main

import (
	"fmt"
	"time"
)

func main() {
	a:=make(map[int]string,10)
	a[1] = "a"
	a[2] = "b"
	a[3] = "c"
	a[4] = "d"
	a[5] = "e"
	a[6] = "f"

	for i:=0 ;i<9;i++ {
		go func() {
			for {
				fmt.Println(a[1])
			}

		}()
	}

	time.Sleep(100*time.Second)
}
