package information_endpoint

import (
	"dbs-api/v1/helpers"

	"github.com/gin-gonic/gin"
)

// ProjectsByTeam GET
// @Summary Get projects by team ProjectsByTeam
// @Schemes
// @Description Get all projects by team requested.
// @Tags Information
// @Accept json
// @Produce json
// @Param        owning_team   query      string  false  "Owning Team as defined in SNOW"
// @Success 200 {string} Example JSON Output
// @Router /information/ProjectsByTeam [get]
func ProjectsByTeam(g *gin.Context) {
	owning_team := g.Request.URL.Query().Get("owning_team")
	db, err := helpers.NewDBClient()
	if err != nil {
		panic(err.Error())
	}

	projects, httpStatus, err := db.GetProjectsByTeam(owning_team)
	if err != nil {
		panic(err.Error())
	}

	g.JSON(httpStatus, projects)
}
