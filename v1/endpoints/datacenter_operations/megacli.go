package endpoint_datacenter_operations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MegaCLI GET
// @Summary Run MegaCLI Commands
// @Schemes
// @Description Run MegaCLI Commands. Only allowable by Database Team and Data Center Team.
// @Tags DataCenter
// @Accept json
// @Produce json
// @Success 200 {string} Example MegaCLI Output
// @Router /datacenter_operations/megacli [get]
func Megacli(g *gin.Context) {
	g.JSON(http.StatusOK, "EXAMPLE MEGACLI_OUTPUT")
}
