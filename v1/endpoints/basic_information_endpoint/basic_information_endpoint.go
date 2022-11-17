package basic_information_endpoint

import (
	clusterDB "dbs-api/v1/helpers/cdb_wrapper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Initialize(big *gin.RouterGroup) {

	big.GET("/GetAllDataCenters", GetAllDataCenters)
	big.GET("/GetAllNetworkZones", GetAllNetworkZones)
	big.GET("/GetAllEnvironments", GetAllEnvironments)
	big.GET("/GetAllFlavors", GetAllFlavors)
	big.GET("/GetAllDatastores", GetAllDatastores)
}

// GetAllDataCenters GET
// @Summary Get all the data centers we can provision to.
// @Schemes
// @Description Get all the data centers we can provision to.
// @Tags Basic Information
// @Accept json
// @Produce json
// @Success 200 {string} Example JSON Output
// @Router /basic_information/GetAllDataCenters [get]
func GetAllDataCenters(g *gin.Context) {
	info, err := clusterDB.GetAllDataCenters()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	}
	g.JSON(http.StatusOK, info)
}

// GetAllNetworkZones GET
// @Summary Get all the network zones we can provision to.
// @Schemes
// @Description Get all the network zones we can provision to.
// @Tags Basic Information
// @Accept json
// @Produce json
// @Success 200 {string} Example JSON Output
// @Router /basic_information/GetAllNetworkZones [get]
func GetAllNetworkZones(g *gin.Context) {
	info, err := clusterDB.GetAllNetworkZones()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	}
	g.JSON(http.StatusOK, info)
}

// GetAllEnvironments GET
// @Summary Get all the environments we can provision to.
// @Schemes
// @Description Get all the environments we can provision to.
// @Tags Basic Information
// @Accept json
// @Produce json
// @Success 200 {string} Example JSON Output
// @Router /basic_information/GetAllEnvironments [get]
func GetAllEnvironments(g *gin.Context) {
	info, err := clusterDB.GetAllEnvironments()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	}
	g.JSON(http.StatusOK, info)
}

// GetAllFlavors GET
// @Summary Get all the flavors we can provision on.
// @Schemes
// @Description Get all the flavors we can provision on.
// @Tags Basic Information
// @Accept json
// @Produce json
// @Success 200 {string} Example JSON Output
// @Router /basic_information/GetAllFlavors [get]
func GetAllFlavors(g *gin.Context) {
	info, err := clusterDB.GetAllFlavors()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	}
	g.JSON(http.StatusOK, info)
}

// GetAllDatastores GET
// @Summary Get all the data stores you can leverage.
// @Schemes
// @Description Get all the data stores you can leverage.
// @Tags Basic Information
// @Accept json
// @Produce json
// @Success 200 {string} Example JSON Output
// @Router /basic_information/GetAllDatastores [get]
func GetAllDatastores(g *gin.Context) {
	info, err := clusterDB.GetAllDatastores()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
	}
	g.JSON(http.StatusOK, info)
}
