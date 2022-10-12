package cluster_endpoint

import (
	"github.com/gin-gonic/gin"
)

func Initialize(cog *gin.RouterGroup) {

	cog.POST("/", Create)
	cog.DELETE("/", Delete)

}
