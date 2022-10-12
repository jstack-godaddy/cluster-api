package information_endpoint

import (
	"github.com/gin-gonic/gin"
)

func Initialize(ig *gin.RouterGroup) {

	ig.GET("/projectsbyteam", ProjectsByTeam)
	ig.GET("/serversbyproject", ServersByProject)

}
