package information_endpoint

import (
	"dbs-api/v1/helpers"
	"fmt"
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
// @Param        owning_team   query      string  false  "Owning Team as defined in SNOW"
// @Success 200 {string} Example JSON Output
// @Router /information/ProjectsByTeam [get]
func ProjectsByTeam(g *gin.Context) {
	var httpStatus int

	owning_team := g.Request.URL.Query().Get("owning_team")
	db, err := helpers.NewDBClient()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	}

	projects, err := db.GetProjectsByTeam(owning_team)
	fmt.Println(projects)
	if err != nil {
		httpStatus = http.StatusBadRequest
		g.String(httpStatus, err.Error())
	} else if len(projects) == 0 {
		httpStatus = http.StatusNotFound
		g.String(httpStatus, "Error: No projects for team provided.")
	}

	g.JSON(httpStatus, projects)
}
