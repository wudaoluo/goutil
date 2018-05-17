package files

import (
	"os"
	"io/ioutil"
	"path/filepath"
)

//返回目录名,文件名
func SplitDirFile(path string) (string, string) {
	return filepath.Dir(path), filepath.Base(path)
}


//判断是否存在
func Exist(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}


//判断是否是文件
func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir()
	}
	return false
}


//判断是否目录
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err == nil {
		return stat.IsDir()
	}
	return false
}


//读取文件字节流
func ReadFileByte(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		//		panic(err)
		return nil, err
	} else {
		defer fi.Close()
		return ioutil.ReadAll(fi)
	}
}