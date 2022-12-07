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

var networkZones = map[string]servers.Network{
	"prd": {UUID: "856447ad-3ce7-4455-90f3-fcfd37f0962c"},
	"mgt": {UUID: "b2e89a04-6c98-4e0c-942d-2edc3f8065a0"},
	"cor": {UUID: "a5ea5932-9159-446e-b6a8-947585d7a044"},
}

var flavors = map[string]string{
	"c8.r16.d200":    "51c496f2-6213-4fec-b893-5ffc5a7f6b4e",
	"c12.r32.d300":   "1e2852bf-70bc-44d0-a85a-05233cca460e",
	"c12.r64.d300":   "e28463d4-3c8d-406d-9a1b-d4ede367bd0e",
	"c16.r96.d900":   "797de406-58d3-48f2-9a4b-77ca04c4b7a0 ",
	"c16.r128.d1200": "5da1a30b-bb12-46f5-83ea-3b6e567b1429 ",
}

var images = map[string]string{
	"alma8": "07d88770-1f3c-4c22-941d-477f853eab89",
	"cent7": "f2df980d-478e-4e08-9513-501c5ee09802",
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
	server, err := client.newServer(networkZone, serverName)
	if err != nil {
		return
	}
	servers = append(servers, server)
	return
}

func (client *NGOSClient) newServer(networkZone string, serverName string) (server *servers.Server, err error) {
	var networks []servers.Network
	networks = append(networks, networkZones[networkZone])
	fmt.Println(networks)
	server, err = servers.Create(client.ServiceClient, servers.CreateOpts{
		Name:      serverName,
		FlavorRef: "a21bcced-5377-49e3-96f8-89e8e9fd5e8d",
		ImageRef:  "07d88770-1f3c-4c22-941d-477f853eab89",
		Networks:  networks,
	}).Extract()
	return
}
