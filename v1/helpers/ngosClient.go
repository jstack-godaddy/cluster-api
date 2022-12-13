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
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/joho/godotenv"
)

type NGOSClient struct {
	ComputeClient *gophercloud.ServiceClient
	NetworkClient *gophercloud.ServiceClient
	ImageService  *gophercloud.ServiceClient
}

var flavors = map[string]string{
	"c8.r16.d200":    "51c496f2-6213-4fec-b893-5ffc5a7f6b4e",
	"c12.r32.d300":   "1e2852bf-70bc-44d0-a85a-05233cca460e",
	"c12.r64.d300":   "e28463d4-3c8d-406d-9a1b-d4ede367bd0e",
	"c16.r96.d900":   "797de406-58d3-48f2-9a4b-77ca04c4b7a0 ",
	"c16.r128.d1200": "5da1a30b-bb12-46f5-83ea-3b6e567b1429 ",
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
	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		return
	}
	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		return
	}
	imageService, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		return
	}

	return &NGOSClient{computeClient, networkClient, imageService}, err
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

	serverNetwork, err := client.getNetworkZones(networkZone)
	if err != nil {
		return
	}
	serverImage, err := client.getLatestImage(os)
	if err != nil {
		return
	}

	for i := 1; i <= quantity; i++ {
		trailing_char := string(foo[i-1])
		serverName := fmt.Sprintf("%s%sl%sdb-%s", dcs[dc], env, shortName, trailing_char)
		server, err := client.newServer(serverName, serverImage, serverNetwork, flavor)
		if err != nil {
			fmt.Printf("Unable to create server: %s\n", err)
			break
		}
		server.Name = serverName
		serversBuilt = append(serversBuilt, server)
		if i == 1 {
			floaterPoolName := fmt.Sprintf("floating-%s", networkZone)
			floaterPool, err := client.getNetworkZones(floaterPoolName)
			if err != nil {
				break
			}
			floater, err = client.reserveFloater(shortName, floaterPool)
			if err != nil {
				break
			}
			go client.associateFloater(server, floater)
		}
	}

	return
}

func (client *NGOSClient) newServer(serverName string, image string, network networks.Network, flavor string) (server *servers.Server, err error) {
	var networks []servers.Network
	networks = append(networks, servers.Network{UUID: network.ID})
	server, err = servers.Create(client.ComputeClient, servers.CreateOpts{
		Name:      serverName,
		FlavorRef: flavors[flavor],
		ImageRef:  image,
		Networks:  networks,
	}).Extract()

	return
}

func (client *NGOSClient) reserveFloater(shortName string, network networks.Network) (fip *floatingips.FloatingIP, err error) {
	createOpts := floatingips.CreateOpts{
		Pool: network.ID,
	}

	fip, err = floatingips.Create(client.ComputeClient, createOpts).Extract()
	//fmt.Println(fip.IP)
	return
}

func (client *NGOSClient) associateFloater(server *servers.Server, fip *floatingips.FloatingIP) {
	time.Sleep(2 * time.Minute)
	associateOpts := floatingips.AssociateOpts{
		FloatingIP: fip.IP,
	}

	err := floatingips.AssociateInstance(client.ComputeClient, server.ID, associateOpts).ExtractErr()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (client *NGOSClient) getNetworkZones(wantedNetworkZone string) (networkZone networks.Network, err error) {
	opts := networks.ListOpts{Name: wantedNetworkZone}

	allPages, err := networks.List(client.NetworkClient, opts).AllPages()
	if err != nil {
		return
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return
	}

	networkZone = allNetworks[0]

	return
}

func (client *NGOSClient) getLatestImage(wantedImageName string) (imageID string, err error) {
	opts := images.ListOpts{Name: wantedImageName}

	allPages, err := images.List(client.ImageService, opts).AllPages()
	if err != nil {
		return
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return
	}

	//fmt.Println(allImages)
	imageID = allImages[0].ID

	return
}
