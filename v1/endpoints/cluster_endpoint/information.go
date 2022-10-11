package cluster_endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	provider, err := ConfigAuth()

	if err == nil {
		g.JSON(http.StatusOK, provider.TokenID)
	} else {
		g.JSON(http.StatusUnauthorized, err)
	}

}
