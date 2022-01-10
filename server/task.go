/**
 * @Author: Resynz
 * @Date: 2022/1/5 15:56
 */
package server

import (
	"github.com/gin-gonic/gin"
	"resync/common"
	"resync/controller/task"
	"resync/enum"
)

func RegisterTaskRoute(route *gin.RouterGroup) {
	route.GET("/list",common.AuthDetection(task.List,enum.MustLogin))
	route.POST("/",common.AuthDetection(task.Add,enum.MustLogin))
	route.PUT("/:id",common.AuthDetection(task.Update,enum.MustLogin))
	route.DELETE("/:id",common.AuthDetection(task.Delete,enum.MustLogin))
	route.POST("/start",common.AuthDetection(task.Start,enum.MustLogin))
	route.POST("/cancel",common.AuthDetection(task.Cancel,enum.MustLogin))
	route.GET("/dump",common.AuthDetection(task.Dump,enum.MustLogin))
}
