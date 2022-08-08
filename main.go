package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() //获取当前执行文件的路径
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe)                 //获取当前执行文件的目录
		filename := uuid.New().String()          //生成一个文件名
		uploads := filepath.Join(dir, "uploads") //拼接upload绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  //创建upload目录
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")                            //拼接文件的绝对路径 不含exe所在目录
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) //将json.raw写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) //返回文件的绝对路径(不含exe所在目录)
	}
}

func main() {
	go func() {
		gin.SetMode(gin.DebugMode) //设置模式
		router := gin.Default()    //新建路由
		//router.GET("/", func(c *gin.Context) { //路由 get方法 c 上下文
		//	c.String(http.StatusOK, "<h1>Hello World</h1>")
		//})
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.POST("/api/v1/texts", TextsController)
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
		router.Run(":8090")
	}()
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://localhost:8090/static/index.html")
	cmd.Start()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	<-chSignal
	cmd.Process.Kill()
	//ui, err := lorca.New("https://www.google.com", "", 800, 600, "--disable-translate")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//select {
	////case <-ui.Done():
	//case <-chSignal:
	//}
	//ui.Close()
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}
