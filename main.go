package main

import (
	"fmt"
	"github.com/beats0/gofileupload/conf"
	"github.com/beats0/gofileupload/models"
	"github.com/beats0/gofileupload/routes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func init() {
	conf.Setup()
	models.Setup()
}

func main() {
	gin.SetMode(conf.ServerSetting.RunMode)

	routersInit := routes.InitRouter()
	readTimeout := conf.ServerSetting.ReadTimeout
	writeTimeout := conf.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", conf.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	log.Printf("[info] start http server listening %s", conf.AppSetting.Url)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("init listen server fail:%v", err)
	}
}
