package cluster_endpoint

import (
	"dbs-api/v1/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CREATE POST
// @Summary Create a server cluster
// @Schemes
// @Description Create a new cluster.
// @Tags Cluster
// @Accept json
// @Produce json
// @Param       dc        		query	string	true 	"Datacenter" Enums(phx, sxb, iad)
// @Param       project   		query	string  true 	"Project Name"
// @Param		name	  		query	string  true    "Cluster Name"
// @Param		shortname		query	string  true    "Abbreviation for naming. No more than 8 characters."
// @Param		networkzone		query	string  true    "Network Zone" Enums(mgt, prd, cor, gcn)
// @Param		privateOrPublic	query	string  true   	"Public or Private floater" Enums(private, public)
// @Success 200 {string} Example JSON Output
// @Router /cluster [post]
func Create(g *gin.Context) {
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

	clusterName := g.Request.URL.Query().Get("name")
	if project == "" {
		g.String(http.StatusBadRequest, "Need to provide name for cluster.")
		return
	}
	fmt.Println(clusterName)

	shortname := g.Request.URL.Query().Get("shortname")
	if project == "" {
		g.String(http.StatusBadRequest, "Need to provide name for cluster.")
		return
	}

	networkzone := g.Request.URL.Query().Get("networkzone")
	if project == "" {
		g.String(http.StatusBadRequest, "Need to provide name for cluster.")
		return
	}

	privateOrPublic := g.Request.URL.Query().Get("privateOrPublic")
	public := false
	if privateOrPublic == "public" {
		public = true
	}

	ngosClient, err := helpers.NewNGOSClient(dc, project)
	if err != nil {
		g.String(http.StatusUnauthorized, err.Error())
		return
	}

	serversCreated, err := ngosClient.NewCluster(shortname, networkzone, public)
	fmt.Println(serversCreated)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	} else {
		g.String(http.StatusOK, "Servers created")
		return
	}
}
