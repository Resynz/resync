/**
 * @Author: Resynz
 * @Date: 2022/1/5 13:56
 */
package server

import (
	"github.com/gin-gonic/gin"
	"resync/common"
	"resync/controller/admin"
	"resync/enum"
)

func RegisterAdminRoute(route *gin.RouterGroup) {
	route.POST("/login", common.AuthDetection(admin.Login))
	route.GET("/info", common.AuthDetection(admin.Info, enum.MustLogin))
	route.GET("/list", common.AuthDetection(admin.List, enum.MustLogin))
	route.POST("/disable", common.AuthDetection(admin.Disable, enum.MustLogin))
	route.POST("/enable", common.AuthDetection(admin.Enable, enum.MustLogin))
	route.POST("/", common.AuthDetection(admin.Add, enum.MustLogin))
	route.PUT("/:id", common.AuthDetection(admin.Update, enum.MustLogin))
	route.DELETE("/:id", common.AuthDetection(admin.Delete, enum.MustLogin))
	route.POST("/modify-passwd", common.AuthDetection(admin.ModifyPasswd, enum.MustLogin))
	route.GET("/log/login", common.AuthDetection(admin.GetAdminLoginLog, enum.MustLogin))
}
