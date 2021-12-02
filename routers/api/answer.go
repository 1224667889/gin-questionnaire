package api

import (
	"github.com/gin-gonic/gin"
	"main/pkg/app"
	"main/pkg/e"
	"net/http"
)

// UploadAnswer 提交问卷
func UploadAnswer(c *gin.Context)  {
	g := app.Gin{C: c}
	params := struct {
		QuestionnaireId string `form:"questionnaire_id" json:"questionnaire_id" xml:"questionnaire_id" binding:"required"`
		Content string `form:"content" json:"content" xml:"content" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}

}