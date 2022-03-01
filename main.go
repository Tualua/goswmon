package main

import (
	"fmt"
	"log"
	// g "github.com/gosnmp/gosnmp"
)

func main() {
	var (
		err error
		// result []g.SnmpPDU
		cfg        *Config
		dhcpLeases DhcpLeases
	)
	/*
		g.Default.Target = "172.18.17.11"
		if err = g.Default.Connect(); err != nil {
			log.Fatalf("Connection error %v", err)
		}
		defer g.Default.Conn.Close()

		oid := "1.3.6.1.2.1.31.1.1.1.1"

		if result, err = g.Default.BulkWalkAll(oid); err != nil {
			log.Fatalf("BulkWalk error %v", err)
		} else {
			for _, v := range result {
				switch v.Type {
				case g.OctetString:
					fmt.Printf("string: %s\n", string(v.Value.([]byte)))
				default:
					fmt.Printf("number: %d\n", g.ToBigInt(v.Value))
				}
			}
		}
	*/
	if cfg, err = NewConfig("config.yaml"); err != nil {
		log.Fatal(err)
	} else {
		ConnectStringROS := fmt.Sprintf("%s:%d", cfg.Sites[1].DhcpServer, cfg.Sites[1].DhcpApiPort)
		if connROS, err := ConnectRos(ConnectStringROS, cfg.Sites[1].DhcpApiLogin, cfg.Sites[1].DhcpApiPassword); err != nil {
			log.Println(err)
		} else {
			defer connROS.Close()
			if dhcpLeases, err = GetLeasesROS(connROS); err != nil {
				log.Println(err)
			} else {
				for _, v := range dhcpLeases {
					if v.Status == "bound" {
						log.Println(v)
					}
				}
			}

		}

	}
}
