package information_endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProjectsByTeam GET
// @Summary Get projects by team ProjectsByTeam
// @Schemes
// @Description Get all projects by team requested.
// @Tags Information
// @Accept json
// @Produce json
// @Param        owning_team   query      string  false  "Datacenter"
// @Success 200 {string} Example JSON Output
// @Router /information/projectsbyteam [get]
func ProjectsByTeam(g *gin.Context) {

	projects := []string{"dbs-infra-dev", "dbs-infra-test", "dbs-infra-prod"}

	g.JSON(http.StatusOK, projects)
}
