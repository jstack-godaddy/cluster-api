package cluster_endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

func Initialize(cog *gin.RouterGroup) {

	cog.GET("/information", Information)

}

func ConfigAuth() (provider *gophercloud.ProviderClient, err error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://phx.openstack.int.gd3p.tools:5000",
		DomainID:         "default",
		Username:         "jstack",
		Password:         "",
		TenantName:       "dbs-infra-dev",
	}

	provider, err = openstack.AuthenticatedClient(opts)

	return
}
