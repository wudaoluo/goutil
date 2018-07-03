package gosys

import (
	"testing"
	"fmt"
)

//str
func Test_GetneiIp_str(t *testing.T) {
	a:=NewAddr(4440)
	a.IntranetAddr()
	fmt.Println(a.GetIPstr())
}

func Test_GetwaiIp_str(t *testing.T) {
	a:=NewAddr(4441)
	a.ExternalAddr()
	fmt.Println(a.GetIPstr())
}

func Test_Getlocal_str(t *testing.T) {
	a:= NewAddr(4002)
	a.LocalAddr()
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaa",a.GetIPstr())
}

//tcp
func Test_GetneiIp_tcp(t *testing.T) {
	a:=NewAddr(4440)
	a.IntranetAddr()
	fmt.Println(a.GetTCPAddr())
}

func Test_GetwaiIp_tcp(t *testing.T) {
	a:=NewAddr(4441)
	a.ExternalAddr()
	fmt.Println(a.GetTCPAddr())
}

//udp
func Test_GetneiIp_udp(t *testing.T) {
	a:=NewAddr(4440)
	a.IntranetAddr()
	fmt.Println(a.GetUDPAddr())
}

func Test_GetwaiIp_udp(t *testing.T) {
	a:=NewAddr(4441)
	a.ExternalAddr()
	fmt.Println(a.GetUDPAddr())
}