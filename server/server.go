/**
 * @Author: Resynz
 * @Date: 2022/1/5 11:55
 */
package server

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"resync/common"
	"resync/config"
	"resync/controller"
	"resync/middleware"
)

//go:embed static
var embedFiles embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embedFiles, "static")
	if err != nil {
		log.Fatalf("load embed file failed! error:%s\n", err.Error())
	}
	return http.FS(fsys)
}

func StartServer() {
	gin.SetMode(config.Conf.Mode)
	app := gin.New()
	app.MaxMultipartMemory = 8 << 20 // 8MB
	_ = app.SetTrustedProxies(nil)

	// 加入中间件
	corConf := cors.DefaultConfig()
	corConf.AllowAllOrigins = true
	corConf.AllowHeaders = []string{
		"Content-Type",
		"token",
	}
	app.Use(cors.New(corConf))

	app.Use(gin.Recovery())

	app.Use(middleware.Auth())

	app.GET("/ping", common.AuthDetection(controller.Ping))

	adminGroup := app.Group("/admin")
	RegisterAdminRoute(adminGroup)

	groupGroup := app.Group("/group")
	RegisterGroupRoute(groupGroup)

	taskGroup := app.Group("/task")
	RegisterTaskRoute(taskGroup)

	logGroup := app.Group("/log")
	RegisterLogRoute(logGroup)

	codeAuthGroup := app.Group("/code_auth")
	RegisterCodeAuthRoute(codeAuthGroup)

	app.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusFound, "/index")
	})
	app.StaticFS("/index", getFileSystem())

	go func() {
		log.Printf("\033[42;30m DONE \033[0m[Resync] Start Success! Port:%d\n", config.Conf.AppPort)
	}()
	if err := app.Run(fmt.Sprintf(":%d", config.Conf.AppPort)); err != nil {
		log.Fatalf("start server failed! error:%s\n", err.Error())
	}
}
