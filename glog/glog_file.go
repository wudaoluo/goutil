// Go support for leveled logs, analogous to https://code.google.com/p/google-glog/
//
// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// File I/O for logs.

package glog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MaxSize is the maximum size of a log file in bytes.


// logDirs lists the candidate directories for new log files.
var logDirs []string

// If non-empty, overrides the choice of directory in which to write logs.
// See createLogDirs for the full list of possible destinations.


func createLogDirs() {
	os.MkdirAll(logDir,0775)  //自动创建目录
	if logDir != "" {
		logDirs = append(logDirs, logDir)
	}
	logDirs = append(logDirs, os.TempDir())
}


var	program  = filepath.Base(os.Args[0])


// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
//func logName(tag string, t time.Time) (name, link string) {
//	name = fmt.Sprintf("%s.%s.%s.log.%s.%04d%02d%02d-%02d%02d%02d.%d",
//		program,
//		host,
//		userName,
//		tag,
//		t.Year(),
//		t.Month(),
//		t.Day(),
//		t.Hour(),
//		t.Minute(),
//		t.Second(),
//		pid)
//	return name, program + "." + tag
//}

func logName(t time.Time) (name,link string) {
	name = fmt.Sprintf("%s.%04d%02d%02d_%02d_%02d_%02d.log",
		program,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
	return name,program + ".log"
}

var onceLogDirs sync.Once

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	onceLogDirs.Do(createLogDirs)
	if len(logDirs) == 0 {
		return nil, "", errors.New("log: no log dirs")
	}
	name, link := logName(t)
	var lastErr error
	for _, dir := range logDirs {
		fname := filepath.Join(dir, name)
		f, err := os.Create(fname)
		if err == nil {
			symlink := filepath.Join(dir, link)
			os.Remove(symlink)        // ignore err
			os.Symlink(name, symlink) // ignore err
			return f, fname, nil
		}
		lastErr = err
	}
	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}
