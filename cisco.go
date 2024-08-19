package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	g "github.com/gosnmp/gosnmp"
)

const (
	OID_CISCO_IFNAME           = "1.3.6.1.2.1.31.1.1.1.1"
	OID_CISCO_VLAN_STATE       = "1.3.6.1.4.1.9.9.46.1.3.1.1.2.1"
	OID_CISCO_IFINDEX          = "1.3.6.1.2.1.17.1.4.1.2"
	OID_CISCO_TPFDBPORT        = "1.3.6.1.2.1.17.4.3.1.2"
	OID_CISCO_ENTPHYSICALENTRY = "1.3.6.1.2.1.47.1.1.1.1.3"
	CISCO_FE                   = ".1.3.6.1.4.1.9.12.3.1.10.48"
	CISCO_GE                   = ".1.3.6.1.4.1.9.12.3.1.10.150"
	CISCO_GE_CONT              = ".1.3.6.1.4.1.9.12.3.1.5.115"
)

type Fdb map[string]string

type CiscoPhysPorts struct {
	FE int
	GE int
}

func GetIfacesNumCisco(ipaddress string, community string) (portconf CiscoPhysPorts, err error) {
	var (
		resSnmp []g.SnmpPDU
	)
	g.Default.Target = ipaddress
	g.Default.Community = community

	if err = g.Default.Connect(); err != nil {
		log.Println(err)
	} else {
		defer g.Default.Conn.Close()
		if resSnmp, err = g.Default.BulkWalkAll(OID_CISCO_ENTPHYSICALENTRY); err != nil {
			log.Println(err)
		} else {
			for _, v := range resSnmp {
				switch fmt.Sprintf("%s", v.Value) {
				case CISCO_FE:
					portconf.FE = portconf.FE + 1
				case CISCO_GE:
					portconf.GE = portconf.GE + 1
				case CISCO_GE_CONT:
					portconf.GE = portconf.GE + 1
				}
			}
		}
	}
	return
}

func GetPortNumFromNameCisco(portname string) (portnum int) {
	s := strings.Split(portname, "/")
	n := s[len(s)-1]
	portnum, _ = strconv.Atoi(n)
	return
}

func GetIfNamesCisco(ipaddress string, community string) (res map[int]string, err error) {
	var (
		resSnmp []g.SnmpPDU
	)
	g.Default.Target = ipaddress
	g.Default.Community = community

	if err = g.Default.Connect(); err != nil {
		log.Println(err)
	} else {
		defer g.Default.Conn.Close()
		if resSnmp, err = g.Default.BulkWalkAll(OID_CISCO_IFNAME); err != nil {
			log.Println(err)
		} else {
			res = make(map[int]string)
			for _, v := range resSnmp {
				ifIndex_oid := strings.Split(v.Name, ".")
				ifIndex, _ := strconv.Atoi(ifIndex_oid[len(ifIndex_oid)-1])
				res[ifIndex] = string(v.Value.([]byte))
			}
		}
	}
	return
}

func GetVlansCisco(ipaddress string, community string) (res []int, err error) {
	var (
		resSnmp []g.SnmpPDU
	)
	g.Default.Target = ipaddress
	g.Default.Community = community

	if err = g.Default.Connect(); err != nil {
		log.Println(err)
	} else {
		defer g.Default.Conn.Close()
		if resSnmp, err = g.Default.BulkWalkAll(OID_CISCO_VLAN_STATE); err != nil {
			log.Println(err)
		} else {
			res = make([]int, 0)
			for _, v := range resSnmp {
				vlan_oid := strings.Split(v.Name, ".")
				vlan, _ := strconv.Atoi(vlan_oid[len(vlan_oid)-1])
				res = append(res, vlan)
			}
		}
	}
	return
}

func asyncGetMappingCisco(s SnmpTask) (res map[int]int) {
	var (
		resSnmp []g.SnmpPDU
		err     error
	)
	snmpParams := &g.GoSNMP{
		Target:    s.IpAddress,
		Port:      uint16(161),
		Community: fmt.Sprintf("%s@%d", s.Community, s.Vlan),
		Version:   g.Version2c,
		Timeout:   time.Duration(6) * time.Second,
	}
	if err = snmpParams.Connect(); err != nil {
		log.Println(err)
	} else {
		defer snmpParams.Conn.Close()
		if resSnmp, err = snmpParams.BulkWalkAll(OID_CISCO_IFINDEX); err != nil {
			log.Println(err)
		} else {
			res = make(map[int]int)
			for _, v := range resSnmp {
				res[GetLastOidOctet(v.Name)] = v.Value.(int)
			}
		}
	}
	return
}

func GetMappingCisco(ipaddress string, community string, vlan int) (res map[int]int, err error) {
	var (
		resSnmp []g.SnmpPDU
	)

	snmpParams := &g.GoSNMP{
		Target:    ipaddress,
		Port:      uint16(161),
		Community: fmt.Sprintf("%s@%d", community, vlan),
		Version:   g.Version2c,
		Timeout:   time.Duration(6) * time.Second,
		//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}
	//g.Default.Target = ipaddress
	//g.Default.Community =
	if err = snmpParams.Connect(); err != nil {
		log.Println(err)
	} else {
		defer snmpParams.Conn.Close()
		if resSnmp, err = snmpParams.BulkWalkAll(OID_CISCO_IFINDEX); err != nil {
			log.Println(err)
		} else {
			res = make(map[int]int)
			for _, v := range resSnmp {
				res[GetLastOidOctet(v.Name)] = v.Value.(int)
			}
		}
	}
	return
}
func asyncGetMappingsCisco(ipaddress string, community string, vlans []int) (res map[int]int) {
	res = make(map[int]int)
	inCh := make(chan SnmpTask)
	go func() {
		defer close(inCh)
		for _, vlan := range vlans {
			inCh <- SnmpTask{
				IpAddress: ipaddress,
				Community: community,
				Vlan:      vlan,
			}
		}
	}()
	outCh := make(chan map[int]int)
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range inCh {
				outCh <- asyncGetMappingCisco(t)
			}
		}()
	}
	go func() {
		defer close(outCh)
		wg.Wait()
	}()
	for s := range outCh {
		for k, v := range s {
			res[k] = v
		}
	}
	return
}

func GetMappingsCisco(ipaddress string, community string, vlans []int) (res map[int]int) {
	res = make(map[int]int)
	for _, vlan := range vlans {
		mapping, _ := GetMappingCisco(ipaddress, community, vlan)
		for k, v := range mapping {
			res[k] = v
		}
	}
	return
}

func GetVlanFdbCisco(s SnmpTask, ifaceMappings map[int]string) (fdb FdbEntries) {
	var (
		resSnmp []g.SnmpPDU
		err     error
	)
	// fdb = make(Fdb)
	snmpParams := &g.GoSNMP{
		Target:    s.IpAddress,
		Port:      uint16(161),
		Community: fmt.Sprintf("%s@%d", s.Community, s.Vlan),
		Version:   g.Version2c,
		Timeout:   time.Duration(6) * time.Second,
	}
	if err = snmpParams.Connect(); err != nil {
		log.Println(err)
	} else {
		defer snmpParams.Conn.Close()
		if resSnmp, err = snmpParams.BulkWalkAll(OID_CISCO_TPFDBPORT); err != nil {
			log.Println(err)
		} else {
			for _, v := range resSnmp {
				fdb.Add(0, ifaceMappings[v.Value.(int)], v.Name, s.Vlan)
			}
		}
	}

	return
}

func GetVlansFdbCisco(ipaddress string, community string, vlans []int, ifaceMappings map[int]string) (fdb FdbEntries) {
	// fdb = make(Fdb)
	inCh := make(chan SnmpTask)
	go func() {
		defer close(inCh)
		for _, vlan := range vlans {
			inCh <- SnmpTask{
				IpAddress: ipaddress,
				Community: community,
				Vlan:      vlan,
			}
		}
	}()
	outCh := make(chan FdbEntries)
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range inCh {
				outCh <- GetVlanFdbCisco(t, ifaceMappings)
			}
		}()
	}
	go func() {
		defer close(outCh)
		wg.Wait()
	}()
	for s := range outCh {
		fdb.Fdb = append(fdb.Fdb, s.Fdb...)
	}
	return
}

func GetFdbCisco(ipaddress string, community string) (res FdbEntries) {
	var (
		vlans []int
		fdb   FdbEntries
	)
	vlans, _ = GetVlansCisco(ipaddress, community)
	ifaces, _ := GetIfNamesCisco(ipaddress, community)
	mappings := asyncGetMappingsCisco(ipaddress, community, vlans)
	ifaceMappings := make(map[int]string)
	for k, v := range mappings {
		ifaceMappings[k] = ifaces[v]
	}
	fdb = GetVlansFdbCisco(ipaddress, community, vlans, ifaceMappings)
	if portConf, err := GetIfacesNumCisco(ipaddress, community); err != nil {
		log.Println(err.Error())
	} else {
		for i := 1; i <= portConf.FE+portConf.GE; i++ {
			res.Fdb = append(res.Fdb, FdbEntry{
				PortNum:    i,
				PortName:   "",
				MacAddress: "Not connected",
				Vlan:       0,
			})
		}
		for i, v := range fdb.Fdb {
			if strings.Contains(v.PortName, "Fa") {
				fdb.Fdb[i].PortNum = GetPortNumFromNameCisco(v.PortName)
			} else {
				fdb.Fdb[i].PortNum = portConf.FE + GetPortNumFromNameCisco(v.PortName)
			}
			res.Fdb[fdb.Fdb[i].PortNum-1] = fdb.Fdb[i]
		}

	}
	return
}
