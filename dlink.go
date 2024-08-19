package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
)

const (
	OID_DOT1QTPFDBENTRY   = ".1.3.6.1.2.1.17.7.1.2.2.1.2"
	OID_DOT1DBASENUMPORTS = ".1.3.6.1.2.1.17.1.2.0"
)

func GetVlanFromOid(oid string) (vlan int) {
	s := strings.Split(oid, ".")
	v := s[len(s)-7]
	vlan, _ = strconv.Atoi(v)
	return
}

func GetNumPortsDlink(ipaddress string, community string) (numports int, err error) {
	var (
		resSnmp *g.SnmpPacket
	)
	snmpParams := &g.GoSNMP{
		Target:    ipaddress,
		Port:      uint16(161),
		Community: community,
		Version:   g.Version2c,
		Timeout:   time.Duration(6) * time.Second,
	}

	if err = snmpParams.Connect(); err != nil {
		log.Println(err)
	} else {
		defer snmpParams.Conn.Close()
		if resSnmp, err = snmpParams.Get([]string{OID_DOT1DBASENUMPORTS}); err != nil {
			log.Println(err)
		} else {
			numports = resSnmp.Variables[0].Value.(int)
		}
	}
	return
}

func GetFdbDlink(ipaddress string, community string) (res FdbEntries) {
	var (
		resSnmp []g.SnmpPDU
		err     error
		fdb     FdbEntries
	)
	snmpParams := &g.GoSNMP{
		Target:    ipaddress,
		Port:      uint16(161),
		Community: community,
		Version:   g.Version2c,
		Timeout:   time.Duration(6) * time.Second,
	}

	if err = snmpParams.Connect(); err != nil {
		log.Println(err)
	} else {
		defer snmpParams.Conn.Close()
		if resSnmp, err = snmpParams.BulkWalkAll(OID_DOT1QTPFDBENTRY); err != nil {
			log.Println(err)
		} else {
			for _, v := range resSnmp {
				if v.Value.(int) > 0 {
					fdb.Add(v.Value.(int), fmt.Sprintf("%d", v.Value), v.Name, GetVlanFromOid(v.Name))
				}
			}
		}
	}
	if NumPorts, err := GetNumPortsDlink(ipaddress, community); err != nil {
		log.Println(err.Error())
	} else {
		for i := 1; i <= NumPorts; i++ {
			res.Fdb = append(res.Fdb, FdbEntry{
				PortNum:    i,
				PortName:   "",
				MacAddress: "Not connected",
				Vlan:       0,
			})
		}
		for i := range fdb.Fdb {

			res.Fdb[fdb.Fdb[i].PortNum-1] = fdb.Fdb[i]
		}
	}

	return
}
