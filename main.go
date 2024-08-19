package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Site struct {
	Name            string `yaml:"name"`
	Prefix          string `yaml:"prefix"`
	Offset          int    `yaml:"offset"`
	DhcpServerType  string `yaml:"dhcp_server_type"`
	DhcpServer      string `yaml:"dhcp_server"`
	DhcpApiPort     int    `yaml:"dhcp_api_port"`
	Community       string `yaml:"community"`
	DhcpApiLogin    string `yaml:"login"`
	DhcpApiPassword string `yaml:"password"`
}

type NbElement struct {
	Name string
	ID   int64
}

type NbList struct {
	Elements []NbElement
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func run(ctx context.Context) (err error) {
	var (
		cfg *Config
	)
	if cfg, err = NewConfig("config.yaml"); err != nil {
		log.Fatal(err)
	} else {
		router := mux.NewRouter().StrictSlash(true)
		router.Use(loggingMiddleware)
		addrString := cfg.Service.Listen + ":" + cfg.Service.Port
		c := NbNewClient(cfg.Netbox.Address, cfg.Netbox.Token)

		router.HandleFunc("/sites", GetSitesList(c))
		router.HandleFunc("/sites/{siteid:[0-9]+}", GetRacksList(c, cfg))
		router.HandleFunc("/rack/{rackid:[0-9]+}", GetDevicesInfo(c, cfg))
		router.HandleFunc("/rack/{rackid:[0-9]+}/csv", GetDevicesInfoCsv(c, cfg))
		router.HandleFunc("/api/minerinfo/", GetMinerInfo)

		srv := &http.Server{
			Addr:    addrString,
			Handler: router,
		}
		go func() {
			if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen:%+s\n", err)
			}
		}()
		log.Printf("server started")
		<-ctx.Done()

		log.Printf("server stopped")
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()

		if err = srv.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("server Shutdown Failed:%+s", err)
		}

		log.Printf("server exited properly")

		if err == http.ErrServerClosed {
			err = nil
		}

	}
	return
}

func (s *Site) GetSiteDhcpLeases() (leases DhcpLeases, err error) {
	leases = make(DhcpLeases)
	switch s.DhcpServerType {
	case "kea":
		var KeaLeases []KeaDhcpLease
		if KeaLeases, err = KeaGetLeasesDhcp4(*s); err != nil {
			log.Println(err.Error())
		} else {
			for _, l := range KeaLeases {
				if l.State == 0 {
					leases[strings.ToUpper(l.HwAddress)] = l.IpAddress
				}
			}
		}
	case "ros":
		var RosLeases []RosDhcpLease
		ConnectStringROS := fmt.Sprintf("%s:%d", s.DhcpServer, s.DhcpApiPort)
		if connROS, err := ConnectRos(ConnectStringROS, s.DhcpApiLogin, s.DhcpApiPassword); err != nil {
			log.Println(err)
		} else {
			defer connROS.Close()
			if RosLeases, err = GetLeasesROS(connROS); err != nil {
				log.Println(err)
			} else {
				for _, l := range RosLeases {
					if l.Status == "bound" {
						leases[strings.ToUpper(l.HwAddress)] = l.IpAddress
					}
				}
			}
		}
	default:
		err = errors.New("unknown DHCP server type")
	}

	return
}

func (s *Site) GetSwitchFdb(swip string, swtype string) (Fdb FdbEntries) {
	switch swtype {
	case "cisco":
		Fdb = GetFdbCisco(swip, s.Community)
	case "d-link":
		Fdb = GetFdbDlink(swip, s.Community)
	}
	sort.Slice(Fdb.Fdb, func(i, j int) bool {
		return Fdb.Fdb[i].PortNum < Fdb.Fdb[j].PortNum
	})
	return
}

func init() {

}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
