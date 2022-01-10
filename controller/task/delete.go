/**
 * @Author: Resynz
 * @Date: 2022/1/6 10:19
 */
package task

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

// Delete 删除任务
func Delete(ctx *common.Context) {
	var task model.Task
	has,err:=db.Handler.XStmt(task.GetTableName()).Where(dbx.Eq("id",ctx.Param("id"))).Get(&task)
	if err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx,code.InvalidRequest,nil)
		return
	}
	if err = db.Handler.Tx(dbx.TxStmts(deleteTask),dbx.TxArg("task_id",task.Id));err!=nil{
		common.HandleResponse(ctx,code.BadRequest,nil,err.Error())
		return
	}
	data:=map[string]bool{
		"result": true,
	}
	common.HandleResponse(ctx,code.SuccessCode,data)
}

func deleteTask(stmt *dbx.TxStmt) error {
	taskId:=stmt.Arg("task_id").(int64)
	var taskLog model.TaskLog
	if err:=stmt.Table(taskLog.GetTableName()).Where(dbx.Eq("task_id",taskId)).Delete(&taskLog);err!=nil{
		return err
	}
	var action model.Action
	if err:=stmt.Table(action.GetTableName()).Where(dbx.Eq("task_id",taskId)).Delete(&action);err!=nil{
		return err
	}
	var taskDetail model.TaskDetail
	if err:=stmt.Table(taskDetail.GetTableName()).Where(dbx.Eq("task_id",taskId)).Delete(&taskDetail);err!=nil{
		return err
	}
	var task model.Task
	if err:=stmt.Table(task.GetTableName()).Where(dbx.Eq("id",taskId)).Delete(&task);err!=nil{
		return err
	}
	return nil
}
