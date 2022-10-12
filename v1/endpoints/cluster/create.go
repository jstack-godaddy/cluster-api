package cluster

import (
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
// @Success 200 {string} Example JSON Output
// @Router /cluster/create [post]
func Create(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
