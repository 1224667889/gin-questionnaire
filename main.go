package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"main/models"
	"main/pkg/setting"
	"main/pkg/util"
	"main/routers"
	"net/http"
)

func init() {
	setting.Setup()
	models.Setup()
	util.Setup()
	//service.InitData()
}


func main() {
	logrus.SetLevel(logrus.DebugLevel)
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logrus.Infoln("start http server listening ", endPoint)

	_ = server.ListenAndServe()
}
