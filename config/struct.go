// Package config /**
package config

type Config struct {
	Mode             string              `json:"mode"`
	AppPort          int                 `json:"app_port"`
	DefaultAdminName string `json:"default_admin_name"`
	DefaultAdminPasswd string `json:"default_admin_passwd"`
	TaskQueueSize int64 `json:"task_queue_size"`
	ActorSize int `json:"actor_size"`
	TaskMap *TaskMap
	Pwd string `json:"pwd"`
}
