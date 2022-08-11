package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"log"
	"net/http"
)

func QrcodsController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal()
		}
		c.Data(http.StatusOK, "image/png", png)
	} else {
		c.Status(http.StatusPreconditionRequired)
	}
}
