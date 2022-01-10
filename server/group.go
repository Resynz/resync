/**
 * @Author: Resynz
 * @Date: 2022/1/5 15:37
 */
package server

import (
	"github.com/gin-gonic/gin"
	"resync/common"
	"resync/controller/group"
	"resync/enum"
)

func RegisterGroupRoute(route *gin.RouterGroup) {
	route.GET("/list",common.AuthDetection(group.List,enum.MustLogin))
	route.POST("/",common.AuthDetection(group.Add,enum.MustLogin))
	route.PUT("/:id",common.AuthDetection(group.Update,enum.MustLogin))
	route.DELETE("/:id",common.AuthDetection(group.Delete,enum.MustLogin))
}
