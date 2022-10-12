package cluster

import (
	"github.com/gin-gonic/gin"
)

func Initialize(cog *gin.RouterGroup) {

	cog.GET("/information", Information)
	cog.POST("/create", Create)

}
