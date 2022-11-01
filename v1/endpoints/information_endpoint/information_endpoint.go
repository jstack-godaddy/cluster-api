package information_endpoint

import (
	"github.com/gin-gonic/gin"
)

func Initialize(ig *gin.RouterGroup) {

	ig.GET("/ProjectsByTeam", ProjectsByTeam)
	ig.GET("/ServersByProject", ServersByProject)
	ig.GET("/GetTeams", GetTeams)

}
