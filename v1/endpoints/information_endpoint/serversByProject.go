package information_endpoint

import (
	"dbs-api/v1/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// ServersByProject GET
// @Summary Get servers by project
// @Schemes
// @Description Displays all servers in a project.
// @Tags Information
// @Accept json
// @Produce json
// @Param       dc        query	string  true 	"Datacenter" Enums(phx, sxb, iad)
// @Param       project   query	string  true 	"Project Name"
// @Success 200 {string} Example JSON Output
// @Router /information/ServersByProject [get]
func ServersByProject(g *gin.Context) {

	dc := g.Request.URL.Query().Get("dc")
	if dc == "" {
		g.JSON(http.StatusBadRequest, "Need to provide datacenter.")
		return
	}

	project := g.Request.URL.Query().Get("project")
	if project == "" {
		g.JSON(http.StatusBadRequest, "Need to provide datacenter.")
		return
	}

	ngosClient, err := helpers.NewNGOSClient(dc, project)

	if err != nil {
		g.JSON(http.StatusUnauthorized, err)
		return
	}

	listOpts := servers.ListOpts{
		AllTenants: false,
	}
	allPages, err := servers.List(ngosClient, listOpts).AllPages()
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
		return
	}
	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
		return
	}

	g.JSON(http.StatusOK, allServers)
}
