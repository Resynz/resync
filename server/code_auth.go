/**
 * @Author: Resynz
 * @Date: 2022/1/7 10:52
 */
package server

import (
	"github.com/gin-gonic/gin"
	"resync/common"
	"resync/controller/code_auth"
	"resync/enum"
)

func RegisterCodeAuthRoute(route *gin.RouterGroup) {
	route.GET("/list",common.AuthDetection(code_auth.List,enum.MustLogin))
	route.POST("/",common.AuthDetection(code_auth.Add,enum.MustLogin))
	route.PUT("/:id",common.AuthDetection(code_auth.Update,enum.MustLogin))
	route.DELETE("/:id",common.AuthDetection(code_auth.Delete,enum.MustLogin))
}
