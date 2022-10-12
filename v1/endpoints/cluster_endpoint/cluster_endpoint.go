package cluster_endpoint

import (
	"github.com/gin-gonic/gin"
)

func Initialize(cg *gin.RouterGroup) {

	cg.POST("/", Create)
	cg.DELETE("/", Delete)

}
