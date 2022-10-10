package endpoint_datacenter_operations

import "github.com/gin-gonic/gin"

func Initialize(dcog *gin.RouterGroup) {

	dcog.GET("/megacli", Megacli)

}
