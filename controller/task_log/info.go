/**
 * @Author: Resynz
 * @Date: 2022/2/8 11:05
 */
package task_log

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/sse"
	"github.com/rosbit/dbx"
	"os"
	"resync/code"
	"resync/common"
	"resync/db"
	"resync/db/model"
)

func GetTaskLogInfo(ctx *common.Context) {
	var taskLog model.TaskLog
	has, err := db.Handler.XStmt(taskLog.GetTableName()).Where(dbx.Eq("id", ctx.Param("id"))).Get(&taskLog)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	if !has {
		common.HandleResponse(ctx, code.InvalidRequest, nil)
		return
	}

	sseRes := sse.Event{Event: "message"}
	sseRes.WriteContentType(ctx.Writer)
	sseRes.Data = map[string]interface{}{
		"action": "open",
	}
	err = sse.Encode(ctx.Writer, sseRes)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	ctx.Writer.Flush()
	// todo 1. 任务状态如果是未知、待执行
	if taskLog.Status == model.TaskStatusUnknown || taskLog.Status == model.TaskStatusPending {
		sseRes.Data = map[string]interface{}{
			"action": "finish",
			"txt":    "该任务无日志",
		}
		err = sse.Encode(ctx.Writer, sseRes)
		if err != nil {
			common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
			return
		}
		ctx.Writer.Flush()
		return
	}

	isQuit := false
	go func() {
		<-ctx.Request.Context().Done()
		isQuit = true
	}()

	logFilePath := fmt.Sprintf("./data/logs/runner/%d/log.log", taskLog.Id)
	lf, err := os.Open(logFilePath)
	if err != nil {
		sseRes.Data = map[string]interface{}{
			"action": "error",
			"msg":    err.Error(),
		}
		err = sse.Encode(ctx.Writer, sseRes)
		if err != nil {
			return
		}
		ctx.Writer.Flush()
		return
	}
	defer lf.Close()
	scanner := bufio.NewScanner(lf)
	for scanner.Scan() && !isQuit {
		sseRes.Data = map[string]interface{}{
			"action": "process",
			"txt":    scanner.Text(),
		}
		err = sse.Encode(ctx.Writer, sseRes)
		if err != nil {
			break
		}
		ctx.Writer.Flush()
	}
	if err = scanner.Err(); err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	// todo 2. 如果任务状态已结束（成功、失败、取消）退出
	if taskLog.Status != model.TaskStatusProcess {
		sseRes.Data = map[string]interface{}{
			"action": "finish",
		}
		err = sse.Encode(ctx.Writer, sseRes)
		if err != nil {
			return
		}
		ctx.Writer.Flush()
		return
	}
	fstat, err := lf.Stat()
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	fsize := fstat.Size()
	// todo 3. 任务正在执行中，需要监听文件变化
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	defer watcher.Close()
	err = watcher.Add(logFilePath)
	if err != nil {
		common.HandleResponse(ctx, code.BadRequest, nil, err.Error())
		return
	}
	for !isQuit {
		event, ok := <-watcher.Events
		if !ok {
			return
		}
		if event.Op == fsnotify.Chmod {
			sseRes.Data = map[string]interface{}{
				"action": "finish",
			}
			err = sse.Encode(ctx.Writer, sseRes)
			if err != nil {
				return
			}
			ctx.Writer.Flush()
			return
		}
		_lf, err := os.Open(logFilePath)
		if err != nil {
			return
		}
		_, err = _lf.Seek(fsize, 0)
		if err != nil {
			return
		}
		scan := bufio.NewScanner(_lf)
		for scan.Scan() && !isQuit {
			fsize += int64(len(scan.Text() + "\n"))
			sseRes.Data = map[string]interface{}{
				"action": "process",
				"txt":    scan.Text(),
			}
			err = sse.Encode(ctx.Writer, sseRes)
			if err != nil {
				break
			}
			ctx.Writer.Flush()
		}
		_lf.Close()
	}
}
