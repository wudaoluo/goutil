package main

func main() {
	var a =  make(chan struct{})
	close(a)
	close(a)
	close(a)
	close(a)

}
