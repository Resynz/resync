// Package config /**
package config

import (
	"os"
	"strconv"
	"sync"
)

var (
	Conf  Config
)

func init() {
	pwd,_:=os.Getwd()
	Conf.Pwd = pwd
	Conf.AppPort = 4930
	Conf.Mode = "debug"
	Conf.DefaultAdminName = "resync"
	Conf.DefaultAdminPasswd = "resynz"
	Conf.TaskQueueSize = 10
	Conf.ActorSize = 2
	Conf.TaskMap = &TaskMap{
		RWMutex: &sync.RWMutex{},
		Map:     make(map[int64]bool),
	}
	if d,err:=strconv.Atoi(os.Getenv("AppPort"));err ==nil && d > 0 {
		Conf.AppPort = d
	}
	if d,err:=strconv.Atoi(os.Getenv("TaskQueueSize"));err ==nil && d > 0 {
		Conf.TaskQueueSize = int64(d)
	}
	if d,err := strconv.Atoi(os.Getenv("ActorSize"));err == nil && d > 0 {
		Conf.ActorSize = d
	}
	if d:=os.Getenv("Mode");d != "" {
		Conf.Mode = d
	}
	if d:=os.Getenv("DefaultAdminName");d != "" {
		Conf.DefaultAdminName = d
	}
	if d:=os.Getenv("DefaultAdminPasswd");d != "" {
		Conf.DefaultAdminPasswd = d
	}
}
