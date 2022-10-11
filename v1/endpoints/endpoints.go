package endpoints

import (
	cluster_endpoint "dbs-api/v1/endpoints/cluster_endpoint"
	endpoint_example "dbs-api/v1/endpoints/endpoint_example"

	"github.com/gin-gonic/gin"
)

func Initialize(v1 *gin.RouterGroup) {

	eg := v1.Group("/example")
	{
		endpoint_example.Initialize(eg)
	}

	cog := v1.Group("/cluster")
	{
		cluster_endpoint.Initialize(cog)
	}

}
