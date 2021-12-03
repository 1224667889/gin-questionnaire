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
	// 搜索问卷
	questionnaire := models.Questionnaire{Id: params.QuestionnaireId}
	if err := models.FindByKey(&questionnaire); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	if questionnaire.AccountId != id {
		g.Response(http.StatusOK, e.FORBIDDEN, "不能操作他人的问卷")
		return
	}
	// 未发布过才能修改
	if questionnaire.HasReleased != "false" {
		g.Response(http.StatusOK, e.FORBIDDEN, "问卷已发布，不允许修改")
		return
	}
	question := models.Question {
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

// AddOption 添加选项
func AddOption(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		Content    string `form:"content" json:"content" xml:"content" binding:"required"`
		QuestionId string `form:"question_id" json:"question_id" xml:"question_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	id := c.MustGet("id").(string)
	// 搜索问题
	question := models.Question{Id: params.QuestionId}
	if err := models.FindByKey(&question); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	// 搜索问卷
	questionnaire := models.Questionnaire{Id: question.QuestionnaireId}
	if err := models.FindByKey(&questionnaire); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	if questionnaire.AccountId != id {
		g.Response(http.StatusOK, e.FORBIDDEN, "不能操作他人的问卷")
		return
	}
	// 未发布过才能修改
	if questionnaire.HasReleased != "false" {
		g.Response(http.StatusOK, e.FORBIDDEN, "问卷已发布，不允许修改")
		return
	}
	option := models.Option{
		Id:         uuid.NewV4().String(),
		Content:    params.Content,
		QuestionId: question.Id,
	}
	if err := models.Insert(&option); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "选项创建失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, option)
	return
}

// DeleteQuestion 删除问题
func DeleteQuestion(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		QuestionId string `form:"question_id" json:"question_id" xml:"question_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	id := c.MustGet("id").(string)
	// 搜索问题
	question := models.Question{Id: params.QuestionId}
	if err := models.FindByKey(&question); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	// 外键搜索问卷
	questionnaire := models.Questionnaire{Id: question.QuestionnaireId}
	if err := models.FindByKey(&questionnaire); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	if questionnaire.AccountId != id {
		g.Response(http.StatusOK, e.FORBIDDEN, "不能操作他人的问卷")
		return
	}
	// 未发布过才能修改
	if questionnaire.HasReleased != "false" {
		g.Response(http.StatusOK, e.FORBIDDEN, "问卷已发布，不允许修改")
		return
	}
	if err := models.Delete(&question); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问题删除失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, question)
	return
}
