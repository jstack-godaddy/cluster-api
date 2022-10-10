package endpoints

import (
	endpoint_datacenter_operations "dbs-api/v1/endpoints/datacenter_operations"
	endpoint_example "dbs-api/v1/endpoints/example"

	"github.com/gin-gonic/gin"
)

func Initialize(v1 *gin.RouterGroup) {

	eg := v1.Group("/example")
	{
		endpoint_example.Initialize(eg)
	}

	dcog := v1.Group("/datacenter_operations")
	{
		endpoint_datacenter_operations.Initialize(dcog)
	}

}
