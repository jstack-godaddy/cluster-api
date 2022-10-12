package information_endpoint

import (
	"github.com/gin-gonic/gin"
)

func Initialize(iog *gin.RouterGroup) {

	iog.GET("/information/ServersByProject", ServersByProject)
	iog.GET("/information/ProjectsByTeam", ProjectsByTeam)

}
