package cluster

import (
	ngosClient "dbs-api/v1/classes"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// Information GET
// @Summary Get information on a cluster
// @Schemes
// @Description Get information on a cluster. Displays all metadata for specified clusters.
// @Tags Cluster
// @Accept json
// @Produce json
// @Success 200 {string} Example JSON Output
// @Router /cluster/information [get]
func Information(g *gin.Context) {
	ngosClient, err := ngosClient.NewClient()

	if err != nil {
		g.JSON(http.StatusUnauthorized, err)
	}

	listOpts := servers.ListOpts{
		AllTenants: false,
	}
	allPages, err := servers.List(ngosClient, listOpts).AllPages()
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
	}
	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
	}

	g.JSON(http.StatusOK, allServers)
}
