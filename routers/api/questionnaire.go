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

// CreateQuestionnaire 创建一个空问卷
func CreateQuestionnaire(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		Title       string `form:"title" json:"title" xml:"title" binding:"required"`
		Subtitle    string `form:"subtitle" json:"subtitle" xml:"subtitle"`
		Description string `form:"description" json:"description" xml:"description"`
		Deadline    string `form:"deadline" json:"deadline" xml:"deadline" binding:"required"`
		NeedLogin   string `form:"need_login" json:"need_login" xml:"need_login" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	wj := models.Questionnaire {
		Id:          uuid.NewV4().String(),
		Title:       params.Title,
		Subtitle:    params.Subtitle,
		Description: params.Description,
		Deadline:    params.Deadline,
		NeedLogin:   params.NeedLogin,
		AccountId:   c.MustGet("id").(string),
	}
	if err := models.Insert(&wj); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问卷创建失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, wj)
	return
}

// AddQuestion 添加问题
func AddQuestion(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		Title           string `form:"title" json:"title" xml:"title" binding:"required"`
		Description     string `form:"description" json:"description" xml:"description"`
		QuestionType    string `form:"question_type" json:"question_type" xml:"question_type"`
		IsRequired      string `form:"is_required" json:"is_required" xml:"is_required" binding:"required"`
		QuestionnaireId string `form:"questionnaire_id" json:"questionnaire_id" xml:"questionnaire_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	id := c.MustGet("id").(string)
	questionnaire := models.Questionnaire{Id: params.QuestionnaireId}
	if err := models.FindByKey(&questionnaire); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	if questionnaire.AccountId != id {
		g.Response(http.StatusOK, e.FORBIDDEN, "不能操作他人的问卷")
		return
	}
	question := models.Question{
		Id:              uuid.NewV4().String(),
		QuestionType:    params.QuestionType,
		Title:           params.Title,
		Description:     params.Description,
		IsRequired:      params.IsRequired,
		QuestionnaireId: questionnaire.Id,
	}
	if err := models.Insert(&question); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问题创建失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, question)
	return
}
