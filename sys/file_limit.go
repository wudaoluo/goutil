package sys

import (
	"syscall"
)


//FileLimit 设置程序的文件打开数
func FileLimit(curLimit,maxLimit uint64) error {
	var rlim syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		return err
	}
	rlim.Cur = curLimit
	rlim.Max = maxLimit
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		return err
	}
	return nil
}
