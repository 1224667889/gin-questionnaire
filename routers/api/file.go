package api

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"main/models"
	"main/pkg/app"
	"main/pkg/e"
	"net/http"
)

// UploadFile 上传文件
func UploadFile(c *gin.Context){
	g := app.Gin{C: c}
	f, err := c.FormFile("file")
	if err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.INVALID_PARAMS, "文件错误")
		return
	}
	id := uuid.NewV4().String()
	if err := c.SaveUploadedFile(f, "./upload/" + id); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.INVALID_PARAMS, "文件错误")
		return
	}
	file := models.File{
		Id:   id,
		Name: f.Filename,
	}
	if err := models.Insert(&file); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "文件保存失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, file)
	return
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context){
	g := app.Gin{C: c}
	uuid := c.DefaultQuery("uuid", "")
	if uuid == "" {
		g.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	id := c.MustGet("id").(string)
	user := models.Account{Id: id}
	if err := models.FindByKey(&user); err != nil {
		g.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "用户信息未找到")
		return
	}
	file := models.File{Id: uuid}
	if err := models.FindByKey(&file); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "uuid解析错误")
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+file.Name)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File("./upload/" + file.Id)
	return
}