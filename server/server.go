package server

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/sharabao13/sync/server/config"
	"github.com/sharabao13/sync/server/controllers"
	"github.com/sharabao13/sync/server/ws"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run() {
	gin.SetMode(gin.DebugMode) //设置模式
	router := gin.Default()    //新建路由
	//router.GET("/", func(c *gin.Context) { //路由 get方法 c 上下文
	//	c.String(http.StatusOK, "<h1>Hello World</h1>")
	//})
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	hub := ws.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.StaticFS("/static", http.FS(staticFiles))
	router.POST("/api/v1/texts", controllers.TextsController)
	router.POST("/api/v1/files", controllers.FilesController)
	router.GET("/api/v1/addresses", controllers.AddressesController)
	router.GET("/uploads/:path", controllers.UploadsController)
	router.GET("/api/v1/qrcodes", controllers.QrcodsController)
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":" + config.GetPort())
}
