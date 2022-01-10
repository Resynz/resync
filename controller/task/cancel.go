/**
 * @Author: Resynz
 * @Date: 2022/1/6 13:37
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
	"time"
)

func Cancel(ctx *common.Context) {
	type formValidate struct {
		Id int64 `form:"id" binding:"required" json:"id"`
	}
	var form formValidate
	if err:=ctx.ShouldBind(&form);err!=nil{
		common.HandleResponse(ctx,code.InvalidParams,nil)
		return
	}
	var taskLog model.TaskLog
	has,err:=db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("id",form.Id)).Get(&taskLog)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.InvalidRequest,nil,"该任务不存在")
		return
	}
	if taskLog.Status > model.TaskStatusProcess {
		common.HandleResponse(ctx,code.InvalidRequest,nil,"该任务不可取消")
		return
	}
	// 1. 将该任务从 TaskMap中 移除
	config.Conf.TaskMap.Delete(taskLog.TaskId)

	if taskLog.Status == model.TaskStatusProcess {
		// cancel
		r,ok:=queue.RunnerMap[taskLog.Id]
		if ok {
			r.Cancel()
		}
	}

	// 2. 修改状态为已取消
	taskLog.Status = model.TaskStatusCancel
	taskLog.EndTime = time.Now().Unix()
	_ , err = db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("id",taskLog.Id)).Cols("status","end_time").Update(&taskLog)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	// todo 3. broadcast

	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}
