package endpoints

import (
	"dbs-api/v1/endpoints/basic_information_endpoint"
	"dbs-api/v1/endpoints/cluster_endpoint"
	"dbs-api/v1/endpoints/endpoint_example"
	"dbs-api/v1/endpoints/project_information_endpoint"
	"dbs-api/v1/endpoints/user_information_endpoint"

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

	uig := v1.Group("/user_information")
	{
		user_information_endpoint.Initialize(uig)
	}

	pig := v1.Group("/project_information")
	{
		project_information_endpoint.Initialize(pig)
	}

	big := v1.Group("/basic_information")
	{
		basic_information_endpoint.Initialize(big)
	}

}
