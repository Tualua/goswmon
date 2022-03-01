package main

import (
	"crypto/tls"
	"log"
	"strings"

	"github.com/go-routeros/routeros"
)

type DhcpLeases map[string]DhcpLease

type DhcpLease struct {
	IPAddress string
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

func GetLeasesROS(ros *routeros.Client) (res DhcpLeases, err error) {
	var (
		rosCommand = "/ip/dhcp-server/lease/print"
		lease      map[string]string
	)
	res = make(DhcpLeases)
	if r, err := ros.RunArgs(strings.Split(rosCommand, " ")); err != nil {
		log.Println(err)
	} else {
		for _, v := range r.Re {
			lease = make(map[string]string)
			for _, j := range v.List {
				lease[j.Key] = j.Value
			}
			res[lease["mac-address"]] = DhcpLease{IPAddress: lease["address"], Hostname: lease["host-name"], Status: lease["status"], LastSeen: lease["last-seen"]}
		}
	}

	return
}
