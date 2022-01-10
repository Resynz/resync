/**
 * @Author: Resynz
 * @Date: 2022/1/7 17:27
 */
package task

import (
	"github.com/gin-contrib/sse"
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/config"
	"resync/db"
	"resync/db/model"
	"resync/queue"
	"time"
)
func Dump(ctx *common.Context) {
	sseRes:=sse.Event{Event: "message"}
	sseRes.WriteContentType(ctx.Writer)
	sseRes.Data = map[string]interface{}{
		"action": "open",
	}
	err:= sse.Encode(ctx.Writer,sseRes)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	ctx.Writer.Flush()
	isQuit:=false
	go func() {
		<- ctx.Request.Context().Done()
		isQuit = true
	}()
	for !isQuit{
		sseRes.Data = makeDumpData()
		err = sse.Encode(ctx.Writer,sseRes)
		if err!=nil{
			break
		}
		ctx.Writer.Flush()
		time.Sleep(time.Second)
	}
}

func makeDumpData() map[string]interface{} {
	type taskObj struct {
		Id int64 `json:"id"`
		Name string `json:"name"`
		PercentTotal int64 `json:"percent_total"`
		Percent int64 `json:"percent"`
	}
	pendingList:=make([]*taskObj,0)

	for k,_:=range config.Conf.TaskMap.Map {
		var task model.Task
		has,err:=db.Handler.XStmt(task.GetTableName()).Where(dbx.Eq("id",k)).Get(&task)
		if err!=nil{
			continue
		}
		if !has {
			continue
		}

		var taskLog model.TaskLog
		has,err = db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("task_id",k)).Desc("id").Get(&taskLog)
		if err!=nil{
			continue
		}
		if !has {
			continue
		}
		l:=&taskObj{
			Id:           taskLog.Id,
			Name:         task.Name,
			PercentTotal: 0,
			Percent:      0,
		}
		pendingList = append(pendingList,l)
	}


	processList:=make([]*taskObj,0)
	for _,v:=range queue.RunnerMap {
		l:=&taskObj{
			Id:           v.Id,
			Name:         v.Task.Name,
			PercentTotal: v.PercentTotal,
			Percent:      v.Percent,
		}
		processList = append(processList,l)
	}
	data:=map[string]interface{}{
		"pending_list": pendingList,
		"process_list": processList,
	}
	return data
}
