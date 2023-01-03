package cluster_endpoint

import (
	"dbs-api/v1/helpers"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type combinedClusterServer struct {
	cluster helpers.Cluster
	servers []helpers.Server
}

// CREATE POST
// @Summary Create a server cluster
// @Schemes
// @Description Create a new cluster.
// @Tags Cluster
// @Accept json
// @Produce json
// @Param       dc        		query	string	true 	"Datacenter the project is in." Enums(phx, sxb, iad)
// @Param       project   		query	string  true 	"Project Name to create the cluster into."
// @Param		name	  		query	string  true    "The friendly name for your identification."
// @Param		shortname		query	string  true    "Abbreviation for naming. Between 3 and 7 characters long." minlength(3) maxlength(7)
// @Param		flavor			query	string	true	"How big do you want it?" Enums(c8.r16.d200,c12.r32.d300,c12.r64.d300,c16.r96.d900,c16.r128.d1200)
// @Param		networkzone		query	string  true    "Network Zone cluster will live in." Enums(prd, prd-public, mgt, cor)
// @Param		os				query	string  true    "Operating System for the cluster." Enums(almalinux8)
// @Param		env				query	string  true    "Environment of cluster. Dev/Test/Stg/OTE/Prod" Enums(d, t, s , o, p)
// @Param		db				query	string  true    "Database Technology being leveraged." Enums(mysql80, mysql57)
// @Success 200 {string} Example JSON Output
// @Router /cluster [post]
func Create(g *gin.Context) {
	dc := g.Request.URL.Query().Get("dc")
	flavor := g.Request.URL.Query().Get("flavor")
	shortname := g.Request.URL.Query().Get("shortname")
	networkzone := g.Request.URL.Query().Get("networkzone")
	os := g.Request.URL.Query().Get("os")
	db := g.Request.URL.Query().Get("db")
	env := g.Request.URL.Query().Get("env")

	project := g.Request.URL.Query().Get("project")
	if project == "" {
		g.String(http.StatusBadRequest, "Need to provide datacenter.")
		return
	}

	clusterName := g.Request.URL.Query().Get("name")
	if clusterName == "" {
		g.String(http.StatusBadRequest, "Need to provide name for cluster.")
		return
	}

	ngosClient, err := helpers.NewNGOSClient(dc, project)
	if err != nil {
		g.String(http.StatusUnauthorized, err.Error())
		return
	}

	serversCreated, fip, err := ngosClient.NewCluster(shortname, dc, db, os, networkzone, env, flavor)
	fmt.Println(serversCreated)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	cdb := helpers.NewClusterDBConn()
	err = cdb.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	clusterReturn, serversReturn, err := cdb.InsertMetadata(clusterName, project, os, networkzone, db, env, flavor, fip, serversCreated)
	if err != nil {
		log.Fatalln(err)
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	jr := combinedClusterServer{
		cluster: clusterReturn,
		servers: serversReturn,
	}
	g.JSON(http.StatusOK, jr)
}
