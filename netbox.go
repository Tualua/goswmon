package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
)

type NbRackInfo struct {
	SiteName     string
	RackName     string
	IpAddress    string
	Manufacturer string
	DeviceType   string
}

func NbNewClient(address string, token string) (c *client.NetBoxAPI) {
	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	transport := httptransport.NewWithClient(address, client.DefaultBasePath, []string{"https"}, httpClient)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+token)
	c = client.New(transport, nil)

	return
}

func NbGetSitesList(c *client.NetBoxAPI) (sites *dcim.DcimSitesListOKBody, err error) {
	reqSites := dcim.NewDcimSitesListParams()
	if resSites, err := c.Dcim.DcimSitesList(reqSites, nil); err != nil {
		log.Println("cannot get sites list from netbox")
	} else {
		sites = resSites.Payload
	}
	return
}

func NbGetRacksList(c *client.NetBoxAPI, siteID string, rackRoleID string) (racks *dcim.DcimRacksListOKBody, err error) {
	rackStatus := "active"
	reqRacks := dcim.NewDcimRacksListParams()
	reqRacks.SetSiteID(&siteID)
	reqRacks.SetRoleID(&rackRoleID)
	reqRacks.SetStatus(&rackStatus)
	if resRacks, err := c.Dcim.DcimRacksList(reqRacks, nil); err != nil {
		log.Printf("cannot get racks list: %v", err)
	} else {
		racks = resRacks.Payload
	}
	return
}

func NbGetRackInfo(c *client.NetBoxAPI, rackID string, switchRoleID string) (rackinfo NbRackInfo, err error) {
	reqDevices := dcim.NewDcimDevicesListParams()
	reqDevices.SetRackID(&rackID)

	if resDevices, err := c.Dcim.DcimDevicesList(reqDevices, nil); err != nil {
		log.Printf("cannot get devices list: %v", err)
	} else {
		rackinfo.SiteName = resDevices.Payload.Results[0].Site.Display
		rackinfo.RackName = resDevices.Payload.Results[0].Rack.Display
		rackinfo.IpAddress = strings.Split(resDevices.Payload.Results[0].PrimaryIP.Display, "/")[0]
		rackinfo.Manufacturer = resDevices.Payload.Results[0].DeviceType.Manufacturer.Display
		rackinfo.DeviceType = resDevices.Payload.Results[0].DeviceType.Display
	}
	return
}
