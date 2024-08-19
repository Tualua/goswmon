package main

import (
	"errors"
	"fmt"
	"log"
	"net/netip"
	"strconv"
	"strings"
)

type SnmpTask struct {
	IpAddress string
	Community string
	Vlan      int
}

type FdbEntry struct {
	PortNum    int
	PortName   string
	MacAddress string
	Vlan       int
}

type FdbEntries struct {
	Fdb []FdbEntry
}

type DhcpLeases map[string]string

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func NewFdbEntry(portnum int, portname string, mac_oid string, vlan int) *FdbEntry {
	return &FdbEntry{
		PortNum:    portnum,
		PortName:   portname,
		MacAddress: GetMacFromOid(mac_oid),
		Vlan:       vlan,
	}
}

func (f *FdbEntries) Add(portnum int, portname string, mac_oid string, vlan int) {
	e := NewFdbEntry(portnum, portname, mac_oid, vlan)
	f.Fdb = append(f.Fdb, *e)
}

/*type DhcpLease struct {
	ClientId  string `json:"client-id"`
	Cltt      int64  `json:"cltt"`
	FqdnFwd   bool   `json:"fqdn-fwd"`
	FqdnRev   bool   `json:"fqdn-rev"`
	Hostname  string `json:"hostname"`
	HwAddress string `json:"hw-address"`
	IpAddress string `json:"ip-address"`
	State     int    `json:"state"`
	SubnetId  int    `json:"subnet-id"`
	ValidLft  int    `json:"valid-lft"`
}*/

func GetLastOidOctet(oid string) (res int) {
	lastString := strings.Split(oid, ".")
	res, _ = strconv.Atoi(lastString[len(lastString)-1])
	return
}

func GetMacFromOid(oid string) (res string) {
	oidSplit := strings.Split(oid, ".")
	macStr := oidSplit[len(oidSplit)-6:]
	for _, v := range macStr {
		o, _ := strconv.Atoi(v)
		res = fmt.Sprintf("%s:%02X", res, o)
	}
	return res[1:]
}

func GetSiteFromIp(cfg *Config, ip string) (site int, err error) {
	site = -1
	if swIp, err := netip.ParseAddr(ip); err != nil {
		log.Println(err.Error())
	} else {
		for i, s := range cfg.Sites {
			if network, err := netip.ParsePrefix(s.Prefix); err != nil {
				log.Println(err.Error())
			} else {
				if network.Contains(swIp) {
					site = i
					break
				}
			}
		}
	}
	if site == -1 {
		err = errors.New("switch IP does not belong to any site")
	}
	return
}
