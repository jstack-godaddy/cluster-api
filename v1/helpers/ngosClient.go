package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/joho/godotenv"
)

type NGOSClient struct {
	ServiceClient *gophercloud.ServiceClient
}

var networkZones = map[string]map[string]servers.Network{
	"phx": {
		"prd":        {UUID: "856447ad-3ce7-4455-90f3-fcfd37f0962c"},
		"prd-public": {UUID: "6617bf94-6201-426a-b0f5-92eb3e6145ee"},
		"mgt":        {UUID: "b2e89a04-6c98-4e0c-942d-2edc3f8065a0"},
		"cor":        {UUID: "a5ea5932-9159-446e-b6a8-947585d7a044"},
	},
}

var floaterPools = map[string]map[string]servers.Network{
	"phx": {
		"prd":        {UUID: "77e6eaa7-5c90-4bd3-b130-18ed67cb645b"},
		"prd-public": {UUID: "484a8183-5029-4dd2-af25-abdca0faef7a"},
		"mgt":        {UUID: "36912f13-3c59-486f-b0b4-35fec1b8f5db"},
		"cor":        {UUID: "630b85ad-c293-4270-88cd-dce51f3a58c3"},
	},
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

var dcs = map[string]string{
	"phx": "p3",
	"sxb": "sxb1",
	"iad": "a2",
}

//var initScript = map[string]string{
//	"mysql8":  "p3",
//	"mysql57": "sxb1",
//}

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

func (client *NGOSClient) NewCluster(shortName string, dc string, db string, os string, networkZone string, env string, flavor string) (serversBuilt []*servers.Server, floater *floatingips.FloatingIP, err error) {
	quantity := 2
	foo := "abcdefghijklmnopqrstuvqxyz"

	switch db {
	case "mysql8":
		quantity = 2
	case "mysql57":
		quantity = 2
	}

	var serverNetwork []servers.Network
	serverNetwork = append(serverNetwork, networkZones[dc][networkZone])

	var floaterPool []servers.Network
	floaterPool = append(floaterPool, floaterPools[dc][networkZone])
	floater, err = client.reserveFloater(shortName, floaterPool)
	if err != nil {
		return
	}

	for i := 1; i <= quantity; i++ {
		trailing_char := string(foo[i-1])
		serverName := fmt.Sprintf("%s%sl%sdb-%s", dcs[dc], env, shortName, trailing_char)
		server, err := client.newServer(serverName, os, serverNetwork, flavor)
		if err != nil {
			break
		}
		serversBuilt = append(serversBuilt, server)
		if i == 1 {
			go client.associateFloater(server, floater)
		}
	}

	return
}

func (client *NGOSClient) newServer(serverName string, os string, network []servers.Network, flavor string) (server *servers.Server, err error) {

	server, err = servers.Create(client.ServiceClient, servers.CreateOpts{
		Name:      serverName,
		FlavorRef: flavors[flavor],
		ImageRef:  images[os],
		Networks:  network,
	}).Extract()
	return
}

func (client *NGOSClient) reserveFloater(shortName string, network []servers.Network) (fip *floatingips.FloatingIP, err error) {
	createOpts := floatingips.CreateOpts{
		Pool: network[0].UUID,
	}

	fip, err = floatingips.Create(client.ServiceClient, createOpts).Extract()
	fmt.Println(fip.IP)
	return
}

func (client *NGOSClient) associateFloater(server *servers.Server, fip *floatingips.FloatingIP) {
	time.Sleep(2 * time.Minute)
	associateOpts := floatingips.AssociateOpts{
		FloatingIP: fip.IP,
	}

	err := floatingips.AssociateInstance(client.ServiceClient, server.ID, associateOpts).ExtractErr()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
