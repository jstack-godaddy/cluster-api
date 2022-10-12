package ngosClient

import (
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/joho/godotenv"
)

func NewClient(identityEndpoint string, tenantName string) (client *gophercloud.ServiceClient, err error) {
	godotenv.Load("./.os_creds.env")
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://phx.openstack.int.gd3p.tools:5000",
		DomainID:         os.Getenv("OS_DOMAIN_ID"),
		Username:         os.Getenv("OS_SVC_USER"),
		Password:         os.Getenv("OS_SVC_PASS"),
		TenantName:       "dbs-infra-dev",
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return
	}

	client, err = openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	return
}
