/**
 * @Author: Resynz
 * @Date: 2022/1/25 11:13
 */
package task

import (
	"github.com/rosbit/dbx"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

// Info 获取任务信息
func Info(ctx *common.Context) {
	var form struct {
		Id int64 `form:"id" binding:"required" json:"id"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		common.HandleResponse(ctx, code.InvalidParams, nil)
		return
	}
	var task model.Task
	has, err := db.Handler.XStmt(task.GetTableName()).Where(dbx.Eq("id", form.Id)).Get(&task)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}
	var taskDetail model.TaskDetail
	_, err = db.Handler.XStmt(taskDetail.GetTableName()).Where(dbx.Eq("task_id", form.Id)).Get(&taskDetail)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	var action model.Action
	var actionList []*model.Action
	err = db.Handler.XStmt(action.GetTableName()).Where(dbx.Eq("task_id", form.Id)).List(&actionList)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	type infoObj struct {
		Task       *model.Task       `json:"task"`
		Detail     *model.TaskDetail `json:"detail"`
		ActionList []*model.Action   `json:"action_list"`
	}
	data := &infoObj{
		Task:       &task,
		Detail:     &taskDetail,
		ActionList: actionList,
	}
	common.HandleResponse(ctx, code.SuccessCode, data)
}
