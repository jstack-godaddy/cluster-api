package helpers

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/joho/godotenv"
)

func NewNGOSClient(identityEndpoint string, tenantName string) (client *gophercloud.ServiceClient, err error) {
	godotenv.Load("./.os_creds.env")
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: fmt.Sprintf("https://%s.openstack.int.gd3p.tools:5000", identityEndpoint),
		DomainID:         os.Getenv("OS_DOMAIN_ID"),
		Username:         os.Getenv("OS_SVC_USER"),
		Password:         os.Getenv("OS_SVC_PASS"),
		TenantName:       tenantName,
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
