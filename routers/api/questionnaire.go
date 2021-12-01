package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"main/models"
	"main/pkg/app"
	"main/pkg/e"
	"main/pkg/util"
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
	questionnaire := models.Questionnaire {
		Id:          uuid.NewV4().String(),
		Title:       params.Title,
		Subtitle:    params.Subtitle,
		Description: params.Description,
		Deadline:    params.Deadline,
		NeedLogin:   params.NeedLogin,
		AccountId:   c.MustGet("id").(string),
	}
	if err := models.Insert(&questionnaire); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问卷创建失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, questionnaire)
	return
}

// DeleteQuestionnaire 删除一个问卷
func DeleteQuestionnaire(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
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
	// 删除问卷
	if err := models.Delete(&questionnaire); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问卷删除失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, nil)
	return
}

// UpdateQuestionnaire 修改一个问卷
func UpdateQuestionnaire(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		Title       string `form:"title" json:"title" xml:"title"`
		Subtitle    string `form:"subtitle" json:"subtitle" xml:"subtitle"`
		Description string `form:"description" json:"description" xml:"description"`
		Deadline    string `form:"deadline" json:"deadline" xml:"deadline"`
		NeedLogin   string `form:"need_login" json:"need_login" xml:"need_login"`
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
	questionnaire.Title = params.Title
	questionnaire.Subtitle = params.Subtitle
	questionnaire.Description = params.Description
	questionnaire.Deadline = params.Deadline
	questionnaire.NeedLogin = params.NeedLogin
	// 删除问卷
	if err := models.Update(&questionnaire); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问卷修改失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, questionnaire)
	return
}

// GetQuestionnaire 获取一个问卷
func GetQuestionnaire(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
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
	// 获取问题
	var questions []models.Question
	if err := models.Find(&questions, fmt.Sprintf("WHERE questionnaire_id = '%s'", questionnaire.Id)); err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问卷查询失败")
		return
	}
	questionnaireMap := util.StructToMapDemo(questionnaire)
	var questionsList []map[string]interface{}
	for _, v := range questions {
		questionsList = append(questionsList, util.StructToMapDemo(v))
	}
	for _, v := range questionsList {
		var options []models.Option
		err := models.Find(&options, fmt.Sprintf("WHERE question_id = '%s'", v["id"].(string)))
		if err == nil {
			v["options"] = options
		}
	}
	questionnaireMap["questions"] = questionsList
	g.Response(http.StatusOK, e.SUCCESS, questionnaireMap)
	return
}

// GetQuestionnaires 获取一堆问卷
func GetQuestionnaires(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		PageSize int `form:"page_size" json:"page_size" xml:"page_size"`
		PageNum  int `form:"page_sum" json:"page_sum" xml:"page_sum"`
	}{}
	// 设置默认值
	if params.PageNum <= 0 {
		params.PageNum = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	id := c.MustGet("id").(string)
	var questionnaires []models.Questionnaire
	if err := models.Find(&questionnaires,
		fmt.Sprintf("WHERE account_id = '%s' LIMIT %d OFFSET %d;", id, params.PageSize, params.PageSize * (params.PageNum - 1)));
		err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问卷查询失败")
		return
	}
	count, err := models.Count(&models.Questionnaire{}, fmt.Sprintf("WHERE account_id = '%s'", id))
	if err != nil {
		logrus.Infoln(err.Error())
		g.Response(http.StatusOK, e.ERROR_DB, "问卷查询失败")
		return
	}
	totalPage := count / params.PageSize
	if totalPage % params.PageSize != 0 {
		totalPage++
	}
	g.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"questionnaires": questionnaires,
		"total_page": totalPage,
	})
	return
}