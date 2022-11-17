package user_information_endpoint

import (
	"dbs-api/v1/helpers"

	"github.com/gin-gonic/gin"
)

func Initialize(uig *gin.RouterGroup) {

	uig.GET("/GetTeams", GetTeams)

}

// GetTeams GET
// @Summary Get teams as designated by SNOW.
// @Schemes
// @Description Displays all team names attached to a username. Will use currently logged in user by default.
// @Tags User Information
// @Accept json
// @Produce json
// @Param       username        query	string  false 	"Username in DC1"
// @Success 200 {string} Example JSON Output
// @Router /user_information/GetTeams [get]
func GetTeams(g *gin.Context) {
	username := g.Request.URL.Query().Get("username")

	teams, httpStatus, err := helpers.GetTeamsFromSNOW(username)
	if err != nil {
		httpStatus = 500
		g.String(httpStatus, err.Error())
	}

	g.JSON(httpStatus, teams)
}
