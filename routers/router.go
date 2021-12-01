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
	apiV1.Use(jwt.JWT())
	{
		apiV1.PUT("/user/detail", api.UpdateUser)
		apiV1.PUT("/user/password", api.ResetPassword)
		apiV1.POST("/user/register", api.RegisterAPI)

		apiV1.POST("/questionnaire", api.CreateQuestionnaire)
		apiV1.POST("/questionnaire/question", api.AddQuestion)
	}
	return r
}
