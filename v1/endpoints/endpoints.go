package endpoints

import (
	"dbs-api/v1/endpoints/cluster_endpoint"
	"dbs-api/v1/endpoints/endpoint_example"
	"dbs-api/v1/endpoints/information_endpoint"

	"github.com/gin-gonic/gin"
)

func Initialize(v1 *gin.RouterGroup) {

	eg := v1.Group("/example")
	{
		endpoint_example.Initialize(eg)
	}

	cg := v1.Group("/cluster")
	{
		cluster_endpoint.Initialize(cg)
	}

	ig := v1.Group("/information")
	{
		information_endpoint.Initialize(ig)
	}

}
