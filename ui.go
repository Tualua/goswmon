package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/netbox-community/go-netbox/netbox/client"
)

type DeviceInfo struct {
	Active     bool
	Rack       string
	PortNum    int
	PortName   string
	MacAddress string
	IpAddress  string
	Vlan       int
	Model      string
	Status     string
	Pool1      string
	Worker1    string
}

type PageDevices struct {
	Site    string
	Rack    string
	Switch  string
	Devices []DeviceInfo
}

func (di DeviceInfo) ToStrings() (res []string) {
	v := reflect.ValueOf(di)
	for i := 0; i < v.NumField(); i++ {
		res = append(res, fmt.Sprintf("%v", v.Field(i).Interface()))
	}
	return res
}

func (di DeviceInfo) Header() (res []string) {
	v := reflect.ValueOf(di)
	for i := 0; i < v.NumField(); i++ {
		res = append(res, fmt.Sprintf("%v", v.Type().Field(i).Name))
	}
	return res
}

func GetSitesList(c *client.NetBoxAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var SiteList NbList
		if NbSites, err := NbGetSitesList(c); err != nil {
			log.Println(err.Error())
		} else {
			sitesTpl := template.Must(template.ParseFiles("templates/sites.html"))
			sitesTpl.Execute(w, NbSites)
		}
	}
}

func GetRacksList(c *client.NetBoxAPI, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		siteID := vars["siteid"]
		if NbRacks, err := NbGetRacksList(c, siteID, cfg.Netbox.RackRole); err != nil {
			log.Println(err.Error())
		} else {
			racksTpl := template.Must(template.ParseFiles("templates/racks.html"))
			racksTpl.Execute(w, NbRacks)
		}
	}
}

func GetDevicesInfo(c *client.NetBoxAPI, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err         error
			pageDevices PageDevices
			RackInfo    NbRackInfo
			// minerType   string
		)
		vars := mux.Vars(r)
		if RackInfo, err = NbGetRackInfo(c, vars["rackid"], cfg.Netbox.SwitchRole); err != nil {
			log.Println("unable to get rack info")
		} else {
			if siteId, err := cfg.FindByName(strings.ToLower(RackInfo.SiteName)); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				site := cfg.Sites[siteId]
				if Leases, err := site.GetSiteDhcpLeases(); err != nil {
					log.Println(err.Error())
				} else {
					SwType := strings.ToLower(RackInfo.Manufacturer)
					Fdb := site.GetSwitchFdb(RackInfo.IpAddress, SwType)
					for _, v := range Fdb.Fdb {
						di := DeviceInfo{
							Active:     v.MacAddress != "Not connected",
							PortNum:    v.PortNum,
							PortName:   v.PortName,
							MacAddress: v.MacAddress,
							IpAddress:  Leases[v.MacAddress],
							Vlan:       v.Vlan,
						}
						if Leases[v.MacAddress] != "" {
							log.Printf("Probing %s", Leases[v.MacAddress])
							m := new(miner)
							m.Init(Leases[v.MacAddress])
							di.Model = m.Model
							di.Status = MinerStatusDesc[m.Status]
							di.Pool1 = m.Pool1
							di.Worker1 = m.Worker1
							log.Printf("%s %f", Leases[v.MacAddress], m.Hashrate/1000/1000)
						}
						pageDevices.Devices = append(pageDevices.Devices, di)
					}
					pageDevices.Site = RackInfo.SiteName
					pageDevices.Rack = RackInfo.RackName
					pageDevices.Switch = fmt.Sprintf("%s %s", RackInfo.Manufacturer, RackInfo.DeviceType)
					SwitchTpl := template.Must(template.ParseFiles("templates/switch.html"))
					SwitchTpl.Execute(w, pageDevices)
				}
			}
		}
	}

}

func GetDevicesInfoCsv(c *client.NetBoxAPI, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err      error
			devices  [][]string
			RackInfo NbRackInfo
			// minerType   string
		)
		vars := mux.Vars(r)
		if RackInfo, err = NbGetRackInfo(c, vars["rackid"], cfg.Netbox.SwitchRole); err != nil {
			log.Println("unable to get rack info")
		} else {
			if siteId, err := cfg.FindByName(strings.ToLower(RackInfo.SiteName)); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				site := cfg.Sites[siteId]
				if Leases, err := site.GetSiteDhcpLeases(); err != nil {
					log.Println(err.Error())
				} else {
					SwType := strings.ToLower(RackInfo.Manufacturer)
					Fdb := site.GetSwitchFdb(RackInfo.IpAddress, SwType)
					devices = append(devices, DeviceInfo{}.Header()[1:])
					for _, v := range Fdb.Fdb {
						di := DeviceInfo{
							Active:     v.MacAddress != "Not connected",
							PortNum:    v.PortNum,
							PortName:   v.PortName,
							MacAddress: v.MacAddress,
							IpAddress:  Leases[v.MacAddress],
							Vlan:       v.Vlan,
						}
						if Leases[v.MacAddress] != "" {
							log.Printf("Probing %s", Leases[v.MacAddress])
							m := new(miner)
							m.Init(Leases[v.MacAddress])
							di.Model = m.Model
							di.Status = MinerStatusDesc[m.Status]
							di.Pool1 = m.Pool1
							di.Worker1 = m.Worker1
							di.Rack = RackInfo.RackName
							devices = append(devices, di.ToStrings()[1:])
							log.Printf("%s %f", Leases[v.MacAddress], m.Hashrate/1000/1000)
						}
					}
					re := regexp.MustCompile(`[^\d\p{Latin}]`)
					filename := strings.Join([]string{re.ReplaceAllString(RackInfo.SiteName, ""), re.ReplaceAllString(RackInfo.RackName, "")}, "-")
					w.Header().Set("Content-Type", "text/csv")
					w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))
					ww := csv.NewWriter(w)
					if err := ww.WriteAll(devices); err != nil {
						log.Println(err.Error())
					}
				}
			}
		}
	}

}
