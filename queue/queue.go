/**
 * @Author: Resynz
 * @Date: 2022/1/6 13:45
 */
package queue

import (
	"log"
	"os"
	"os/signal"
	"resync/config"
	"resync/db/model"
	"resync/runner"
	"sync"
	"syscall"
)



var (
	TaskQueue chan *model.TaskLog
	ActorList []*Actor
	RunnerMap map[int64]*runner.Runner

)

func InitQueues() {
	RunnerMap = make(map[int64]*runner.Runner)
	RegisterSignalHandler()
	TaskQueue = make(chan *model.TaskLog,config.Conf.TaskQueueSize)
	for i:=1;i<=config.Conf.ActorSize;i++ {
		act:=&Actor{
			Id:       i,
			StopFlag: false,
			StopChan: make(chan bool),
		}
		ActorList = append(ActorList,act)
		go act.Start()
	}
}


func RegisterSignalHandler()  {
	sign:=make(chan os.Signal,1)
	signal.Notify(sign,syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	go func() {
		defer func() {
			log.Println("程序退出")
			os.Exit(0)
		}()
		_ = <-sign
		log.Println("接收到退出信号，准备退出...")
		wg:=sync.WaitGroup{}
		wg.Add(len(ActorList))
		for _,v:=range ActorList {
			v := v
			go func() {
				v.Stop()
				wg.Done()
			}()
		}
		wg.Wait()
	}()
}