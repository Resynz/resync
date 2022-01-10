/**
 * @Author: Resynz
 * @Date: 2022/1/6 10:28
 */
package task

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/config"
	"resync/db"
	"resync/db/model"
	"resync/queue"
)

// Start 启动任务
func Start(ctx *common.Context) {
	type formValidate struct {
		Id int64 `form:"id" binding:"required" json:"id"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	var task model.Task
	has,err:=db.Handler.XStmt(task.GetTableName()).Where(dbx.Eq("id",form.Id)).Get(&task)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.InvalidRequest,nil,"任务不存在或已被删除")
		return
	}
	if config.Conf.TaskMap.Exists(task.Id) {
		common.HandleResponse(ctx,code.InvalidRequest,nil,"该任务正在处理中，请勿重复添加")
		return
	}
	taskLog:=model.TaskLog{
		Id:        0,
		TaskId:    task.Id,
		Status:    model.TaskStatusPending,
		StartTime: 0,
		EndTime:   0,
		CreatorId: ctx.Auth.Id,
	}
	if err = db.Handler.XStmt(taskLog.GetTableName()).Insert(&taskLog);err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	config.Conf.TaskMap.Set(taskLog.TaskId)
	// push to queue
	queue.TaskQueue<-&taskLog
	// todo broadcast sse

	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
