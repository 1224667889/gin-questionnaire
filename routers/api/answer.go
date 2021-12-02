package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"main/models"
	"main/mongo"
	"main/pkg/app"
	"main/pkg/e"
	"main/pkg/util"
	"net/http"
	"strings"
)

// UploadAnswer 提交答卷
func UploadAnswer(c *gin.Context)  {
	g := app.Gin{C: c}
	params := struct {
		QuestionnaireId string        `form:"questionnaire_id" json:"questionnaire_id" xml:"questionnaire_id" binding:"required"`
		Answers         []interface{} `form:"answers" json:"answers" xml:"answers" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, err.Error())
		return
	}
	// 搜索问卷
	questionnaire := models.Questionnaire{Id: params.QuestionnaireId}
	if err := models.FindByKey(&questionnaire); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, "参数错误")
		return
	}
	// 判断问卷开放
	if questionnaire.IsOpen != "false" {
		g.Response(http.StatusOK, e.FORBIDDEN, "问卷未开放")
		return
	}
	var id string
	// 判断是否需要登录填写
	if questionnaire.NeedLogin != "false" {
		token := c.GetHeader("Authorization")
		if token == "" {
			g.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "该问卷需要登录")
			return
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				g.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "该问卷需要登录")
			default:
				g.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "该问卷需要登录")
			}
			return
		}
		id = claims.Id
	}
	// 获取问题格式
	var questions []models.Question
	if err := models.Find(&questions, fmt.Sprintf("WHERE questionnaire_id = '%s'", questionnaire.Id)); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, "问题查询失败")
		return
	}
	var totalRequired int
	questionsMap := make(map[string]models.Question)
	// 映射问题到id，统计必做题总数
	for _, question := range questions {
		// 判断是否该题是当前问卷的
		if question.QuestionnaireId != questionnaire.Id {
			g.Response(http.StatusOK, e.INVALID_PARAMS, "问卷选项无法对应")
			return
		}
		// 映射到map，方便查找
		questionsMap[question.Id] = question
		if question.IsRequired == "true" {
			// 统计必做题数量
			totalRequired++
		}
	}
	for _, v := range params.Answers {
		answer := v.(map[string]interface{})
		questionId := answer["id"].(string)
		content := answer["content"].(string)
		if question, ok := questionsMap[questionId]; ok {
			switch question.QuestionType {
			case "单选", "多选":
				optionIdList := strings.Split(content, ",")
				var options []models.Option
				if err := models.Find(&options, fmt.Sprintf("WHERE question_id = '%s'", questionId)); err != nil {
					logrus.Infoln(err.Error())
					g.Response(http.StatusOK, e.ERROR_DB, "问题选项查询失败")
					return
				}
				if len(optionIdList) != len(options) {
					g.Response(http.StatusOK, e.INVALID_PARAMS, "选项数量异常")
					return
				}
				var sumChoose int
				for _, optionId := range optionIdList {
					switch optionId {
					case "0":
					case "1":
						sumChoose++
					default:
						g.Response(http.StatusOK, e.INVALID_PARAMS, "选项内容异常")
						return
					}
				}
				if sumChoose <= 0 || (sumChoose != 1 && question.QuestionType == "单选") {
					g.Response(http.StatusOK, e.INVALID_PARAMS, "选择选项数量异常")
					return
				}
			case "填空":
				// Todo: 验证内容，暂不做判断
			case "文件":
				file := models.File{Id: content}
				err := models.FindByKey(&file)
				if err != nil {
					g.Response(http.StatusOK, e.ERROR_DB, "文件不存在")
					return
				}
			default:
				g.Response(http.StatusOK, e.INVALID_PARAMS, "选项类型错误")
				return
			}
			if question.IsRequired == "true" {
				totalRequired--
			}
		} else {
			g.Response(http.StatusOK, e.INVALID_PARAMS, "选项无法对应")
			return
		}
	}
	if totalRequired != 0 {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "存在未填写必选")
		return
	}
	var save mongo.SaveAnswers
	save.Id = id
	save.Answers = params.Answers
	collection := mongo.MongoDB.Collection(questionnaire.Id)
	res, err := collection.InsertOne(context.TODO(), save)
	if err != nil {
		logrus.Infoln("collection insert one data failed: ", err)
		g.Response(http.StatusOK, e.ERROR_DB, "插入失败")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, res.InsertedID)
	return
}

// GetAnswers 获取答卷
func GetAnswers(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		QuestionnaireId string        `form:"questionnaire_id" json:"questionnaire_id" xml:"questionnaire_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, err.Error())
		return
	}
	// 搜索问卷
	questionnaire := models.Questionnaire{Id: params.QuestionnaireId}
	if err := models.FindByKey(&questionnaire); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, "参数错误")
		return
	}
	collection := mongo.MongoDB.Collection(questionnaire.Id)
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		logrus.Infoln("collection find data failed: ", err)
		g.Response(http.StatusOK, e.ERROR_DB, "查询失败")
		return
	}
	var saveAnswers []mongo.LoadAnswers
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var answer mongo.LoadAnswers
		err := cur.Decode(&answer)
		if err != nil {
			logrus.Infoln("collection find data failed: ", err)
			g.Response(http.StatusOK, e.ERROR_DB, "解码失败")
			return
		}
		saveAnswers = append(saveAnswers, answer)
	}
	g.Response(http.StatusOK, e.SUCCESS, saveAnswers)
	return
}