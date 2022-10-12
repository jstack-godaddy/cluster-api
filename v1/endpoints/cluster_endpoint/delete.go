package cluster_endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DELETE DELETE
// @Summary Delete a server cluster
// @Schemes
// @Description Delete a new cluster.
// @Tags Cluster
// @Accept json
// @Produce json
// @Success 200 {string} Example JSON Output
// @Router /cluster [delete]
func Delete(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
