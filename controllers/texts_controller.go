package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() //获取当前执行文件的路径
		if err != nil {
			log.Fatal()
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
