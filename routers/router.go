package routers

import (
	"github.com/gin-gonic/gin"
	"main/middleware/jwt"
	"main/routers/api"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")
	apiV1.POST("/login", api.LoginAPI)
	apiV1.POST("/user/register", api.RegisterAPI)
	apiV1.POST("/questionnaire/answer", api.UploadAnswer)
	apiV1.POST("/file", api.UploadFile)
	apiV1.Use(jwt.JWT())
	{
		apiV1.PUT("/user/detail", api.UpdateUser)
		apiV1.PUT("/user/password", api.ResetPassword)

		apiV1.GET("/questionnaire", api.GetQuestionnaire)
		apiV1.GET("/questionnaires", api.GetQuestionnaires)
		apiV1.POST("/questionnaire", api.CreateQuestionnaire)
		apiV1.PUT("/questionnaire", api.UpdateQuestionnaire)
		apiV1.DELETE("/questionnaire", api.DeleteQuestionnaire)
		apiV1.GET("/questionnaire/answer", api.GetAnswers)

		apiV1.PUT("/questionnaire/release", api.ReleaseQuestionnaire)

		apiV1.POST("/question", api.AddQuestion)
		apiV1.DELETE("/question", api.DeleteQuestion)
		apiV1.POST("/question/option", api.AddOption)

		apiV1.GET("/file", api.DownloadFile)
	}
	return r
}
