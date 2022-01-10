/**
 * @Author: Resynz
 * @Date: 2022/1/6 13:56
 */
package queue

import (
	"context"
	"fmt"
	"github.com/rosbit/dbx"
	"log"
	"resync/config"
	"resync/db"
	"resync/db/model"
	"resync/runner"
	"time"
)

type IActor interface {
	Start()
	Stop()
}

type Actor struct {
	Id int
	StopFlag bool
	StopChan chan bool
}

func (s *Actor) log(txt string)  {
	log.Printf("[Actor:%d] %s\n",s.Id,txt)
}

func (s *Actor) Start()  {
	s.log("start!")
	for !s.StopFlag {
		taskLog:=<-TaskQueue
		if taskLog == nil {
			continue
		}
		if !config.Conf.TaskMap.Exists(taskLog.TaskId) {
			continue
		}
		// 将该任务从 TaskMap中 移除
		config.Conf.TaskMap.Delete(taskLog.TaskId)
		var task model.Task
		has,err:=db.Handler.XStmt(task.GetTableName()).Where(dbx.Eq("id",taskLog.TaskId)).Get(&task)
		if err!=nil{
			s.log(fmt.Sprintf("find task failed! err:%s",err.Error()))
			continue
		}
		if !has {
			s.log("task not found!")
			continue
		}
		var taskDetail model.TaskDetail
		has,err = db.Handler.XStmt(taskDetail.GetTableName()).Where(dbx.Eq("task_id",task.Id)).Get(&taskDetail)
		if err!=nil{
			s.log(fmt.Sprintf("find task_detail failed! err:%s",err.Error()))
			continue
		}
		if !has {
			s.log("task_detail not found!")
			continue
		}

		var admin model.Admin
		has,err = db.Handler.XStmt(admin.GetTableName()).Where(dbx.Eq("id",taskLog.CreatorId)).Get(&admin)
		if err!=nil{
			s.log(fmt.Sprintf("find admin failed! err:%s",err.Error()))
			continue
		}
		if !has {
			s.log("admin not found!")
			continue
		}
		var action model.Action
		var actionList []*model.Action
		err = db.Handler.XStmt(action.GetTableName()).Where(dbx.Eq("task_id",taskLog.TaskId)).List(&actionList)
		if err!=nil{
			s.log(fmt.Sprintf("find actions failed! err:%s",err.Error()))
			continue
		}
		taskLog.StartTime = time.Now().Unix()
		taskLog.Status = model.TaskStatusProcess
		if _,err = db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("id",taskLog.Id),dbx.Eq("status",model.TaskStatusPending)).Cols("status","start_time").Update(taskLog);err!=nil{
			continue
		}
		ctx,cancel:=context.WithCancel(context.Background())
		r:=&runner.Runner{
			Id:           taskLog.Id,
			ExecutorId:   taskLog.CreatorId,
			ExecutorName: admin.Name,
			Task:         &task,
			Detail:       &taskDetail,
			ActionList:   actionList,
			Context: ctx,
			Cancel: cancel,
			Percent: 0,
			PercentTotal: 3 + int64(len(actionList)),
		}
		RunnerMap[r.Id] = r
		if err = r.Run();err!=nil{
			s.log(fmt.Sprintf("run failed! err:%s",err.Error()))
			taskLog.EndTime = time.Now().Unix()
			taskLog.Status = model.TaskStatusFailed
			db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("id",taskLog.Id),dbx.Eq("status",model.TaskStatusProcess)).Cols("status","end_time").Update(taskLog)
			delete(RunnerMap,r.Id)
			continue
		}
		taskLog.EndTime = time.Now().Unix()
		taskLog.Status = model.TaskStatusSuccess
		db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("id",taskLog.Id),dbx.Eq("status",model.TaskStatusProcess)).Cols("status","end_time").Update(taskLog)
		delete(RunnerMap,r.Id)
		s.log("done")
	}
	s.StopChan <- true
}

func (s *Actor) Stop()  {
	s.log("stopping ...")
	s.StopFlag = true
	go func() {
		time.Sleep(time.Millisecond * 100)
		TaskQueue<-nil
	}()
	<-s.StopChan
	s.log("stopped.")
}