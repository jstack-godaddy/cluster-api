package v1

import (
	docs "dbs-api/docs"
	endpoints "dbs-api/v1/endpoints"

	gin "github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title DBS Linux Cluster API
// @version 1.0

// @contact.name #mysql in Slack
// @BasePath /api/v1
func Initialize(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	{
		endpoints.Initialize(v1)
	}

	r.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
