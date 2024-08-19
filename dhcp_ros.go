package main

import (
	"crypto/tls"
	"log"
	"strings"

	"github.com/go-routeros/routeros"
)

type RosDhcpLease struct {
	HwAddress string
	IpAddress string
	Hostname  string
	Status    string
	LastSeen  string
}

func ConnectRos(address string, username string, password string) (*routeros.Client, error) {
	var (
		tlsConfig tls.Config
	)
	tlsConfig.InsecureSkipVerify = true
	return routeros.DialTLS(address, username, password, &tlsConfig)
}

func GetLeasesROS(ros *routeros.Client) (res []RosDhcpLease, err error) {
	var (
		rosCommand = "/ip/dhcp-server/lease/print"
		lease      map[string]string
	)
	if r, err := ros.RunArgs(strings.Split(rosCommand, " ")); err != nil {
		log.Println(err)
	} else {
		for _, v := range r.Re {
			lease = make(map[string]string)
			for _, j := range v.List {
				lease[j.Key] = j.Value
			}
			res = append(
				res,
				RosDhcpLease{HwAddress: lease["mac-address"], IpAddress: lease["address"], Hostname: lease["host-name"], Status: lease["status"], LastSeen: lease["last-seen"]})
		}
	}

	return
}
