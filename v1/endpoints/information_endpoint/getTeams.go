package information_endpoint

import (
	"dbs-api/v1/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTeams GET
// @Summary Get teams as designated by SNOW.
// @Schemes
// @Description Displays all team names attached to a username. Will use currently logged in user by default.
// @Tags Information
// @Accept json
// @Produce json
// @Param       username        query	string  false 	"Username in DC1"
// @Success 200 {string} Example JSON Output
// @Router /information/GetTeams [get]
func GetTeams(g *gin.Context) {
	username := g.Request.URL.Query().Get("username")
	
	teams,err := helpers.GetTeams(username)
	if err != nil {
		return
	}
	
	g.JSON(http.StatusOK, teams)
}
