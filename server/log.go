/**
 * @Author: Resynz
 * @Date: 2022/2/8 11:02
 */
package server

import (
	"github.com/gin-gonic/gin"
	"resync/common"
	"resync/controller/task_log"
	"resync/enum"
)

func RegisterLogRoute(route *gin.RouterGroup) {
	route.GET("/info/:id", common.AuthDetection(task_log.GetTaskLogInfo, enum.MustLogin))
	route.GET("/list", common.AuthDetection(task_log.GetTaskLogList, enum.MustLogin))
}
