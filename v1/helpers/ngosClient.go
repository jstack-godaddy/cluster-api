package helpers

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/joho/godotenv"
)

type NGOSClient struct {
	ServiceClient *gophercloud.ServiceClient
}

func NewNGOSClient(identityEndpoint string, tenantName string) (ngosclient *NGOSClient, err error) {
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

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	return &NGOSClient{client}, err
}

func (client *NGOSClient) NewCluster(shortName string, networkZone string, public bool) (servers []*servers.Server, err error) {
	serverName := "p3pltestboxdb-a"
	server, err := client.NewServer(networkZone, serverName)
	if err != nil {
		return
	}
	servers = append(servers, server)
	return
}

func (client *NGOSClient) NewServer(networkZone string, serverName string) (server *servers.Server, err error) {
	networks := []servers.Network{
		{UUID: "856447ad-3ce7-4455-90f3-fcfd37f0962c"},
	}
	fmt.Println(networks)
	server, err = servers.Create(client.ServiceClient, servers.CreateOpts{
		Name:      serverName,
		FlavorRef: "a21bcced-5377-49e3-96f8-89e8e9fd5e8d",
		ImageRef:  "07d88770-1f3c-4c22-941d-477f853eab89",
		Networks:  networks,
	}).Extract()
	return
}
