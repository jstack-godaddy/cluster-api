package project_information_endpoint

import (
	"dbs-api/v1/helpers"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

var Cdb *helpers.ClusterDB

func Initialize(pig *gin.RouterGroup) {

	Cdb = helpers.NewClusterDBConn()
	err := Cdb.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	pig.GET("/ProjectsByTeam", ProjectsByTeam)
	pig.GET("/ServersByProject", ServersByProject)
	pig.GET("/ServersByProjectRaw", ServersByProjectRaw)

}

// ProjectsByTeam GET
// @Summary Get projects by team ProjectsByTeam
// @Schemes
// @Description Get all projects by team requested.
// @Tags Project Information
// @Accept json
// @Produce json
// @Param        owning_team   query      string  false  "Owning Team as defined in SNOW"
// @Success 200 {string} Example JSON Output
// @Router /project_information/ProjectsByTeam [get]
func ProjectsByTeam(g *gin.Context) {
	owning_team := g.Request.URL.Query().Get("owning_team")

	projects, err := Cdb.GetProjectsByTeam(owning_team)
	fmt.Println(projects)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	} else if len(projects) == 0 {
		g.String(http.StatusNotFound, "Error: No projects for team provided.")
	}

	g.JSON(http.StatusOK, projects)
}

// ServersByProject GET
// @Summary Get servers by project.
// @Schemes
// @Description Displays all servers in a project by directly querying our metadata.
// @Tags Project Information
// @Accept json
// @Produce json
// @Param       project_id   query	int  true 	"Project ID Number"
// @Success 200 {string} Example JSON Output
// @Router /project_information/ServersByProject [get]
func ServersByProject(g *gin.Context) {

	project_id := g.Request.URL.Query().Get("project_id")
	if project_id == "" {
		g.String(http.StatusBadRequest, "Need to provide project ID.")
		return
	}

	servers, err := Cdb.GetServersByProject(project_id)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	} else if len(servers) == 0 {
		g.String(http.StatusNotFound, "No servers found in project.")
	}

	g.JSON(http.StatusOK, servers)
}

// ServersByProjectRaw GET
// @Summary Get servers by project.
// @Schemes
// @Description Displays all servers in a project by directly querying Openstack.
// @Tags Project Information
// @Accept json
// @Produce json
// @Param       dc        query	string  true 	"Datacenter" Enums(phx, sxb, iad)
// @Param       project   query	string  true 	"Project Name"
// @Success 200 {string} Example JSON Output
// @Router /project_information/ServersByProjectRaw [get]
func ServersByProjectRaw(g *gin.Context) {

	dc := g.Request.URL.Query().Get("dc")
	if dc == "" {
		g.String(http.StatusBadRequest, "Need to provide datacenter.")
		return
	}

	project := g.Request.URL.Query().Get("project")
	if project == "" {
		g.String(http.StatusBadRequest, "Need to provide datacenter.")
		return
	}

	ngosClient, err := helpers.NewNGOSClient(dc, project)

	if err != nil {
		g.String(http.StatusUnauthorized, err.Error())
		return
	}

	listOpts := servers.ListOpts{
		AllTenants: false,
	}
	allPages, err := servers.List(ngosClient, listOpts).AllPages()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, allServers)
}
