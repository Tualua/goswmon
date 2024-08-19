package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"strings"
	"time"
)

const (
	defaultCgMinerApiPort                     = 4028
	ApiCmd                MinerApiCommandType = 0
	ApiCommand            MinerApiCommandType = 1
)

type MinerStatus int32

const (
	MinerOK            MinerStatus = 0
	MissingHashBlade   MinerStatus = 1
	HaveDeadHashBlades MinerStatus = 2
	Timeout            MinerStatus = 3
	Unknown            MinerStatus = 4
)

var MinerStatusDesc = [...]string{
	"OK",
	"Missing Hash Blade(s)",
	"Have Dead Hash Blade(s)",
	"Timeout",
	"Unknown",
}

type MinerApiCommandType int32

type MinerApiCommand struct {
	Cmd     string `json:"cmd,omitempty"`
	Command string `json:"command,omitempty"`
}

type minerHashblade struct {
	Hashrate    int64
	Temperature int
	Alive       bool
}

type minerFan struct {
	Id    string
	Speed int
}

type miner struct {
	Type           MinerType
	Model          string `json:"model"`
	IpAddress      string
	ApiPort        int
	ApiCommandType MinerApiCommandType
	MacAddress     string
	NumBlades      int
	NumFans        int
	Hashrate       float64 `json:"hashrate"`
	Power          int     `json:"power"`
	Status         MinerStatus
	FwVersion      string
	HashBlades     []minerHashblade
	Fans           []minerFan
	Worker1        string
	Pool1          string
}

type CgMinerStatus struct {
	Status      string `json:"STATUS"`
	Timestamp   int64  `json:"When"`
	Msg         string `json:"Msg"`
	Description string `json:"Description"`
}

type MinerTypeReply struct {
	Status []CgMinerStatus `json:"STATUS"`
}

type MinerType int32

const (
	UnknownMiner MinerType = 0
	Antminer     MinerType = 1
	Whatsminer   MinerType = 2
	Innosilicon  MinerType = 3
	Avalon       MinerType = 4
)

func (c *MinerApiCommand) SetCommand(command string, commandType MinerApiCommandType) {
	switch commandType {
	case ApiCmd:
		c.Command = ""
		c.Cmd = command
	case ApiCommand:
		c.Command = command
		c.Cmd = ""
	}
}

func (c *MinerApiCommand) GetCommand() (res []byte) {
	var (
		err error
	)
	if res, err = json.Marshal(c); err != nil {
		fmt.Println(err.Error())
	}
	return
}

func (m *miner) Init(ipAddress string) {
	var (
		apiCommand MinerApiCommand
		apiReply   MinerTypeReply
		rawReply   []byte
		err        error
	)
	m.IpAddress = ipAddress
	m.Type = UnknownMiner
	m.Model = "Unknown"
	m.Status = Unknown
	//try 4028 port first
	apiCommand.SetCommand("summary", ApiCmd)
	if rawReply, err = SendApiCommand(m.IpAddress, defaultCgMinerApiPort, apiCommand.GetCommand()); err != nil {
		switch cause := err.(type) {
		case net.Error:
			if cause.Timeout() {
				m.Status = Timeout
			}
		default:
			log.Println(err.Error())
		}
	} else {
		json.Unmarshal(rawReply, &apiReply)
		if len(apiReply.Status) > 0 {
			m.ApiPort = defaultCgMinerApiPort
			if strings.Contains(apiReply.Status[0].Msg, "Summary") {
				m.Type = Whatsminer
				m.ApiCommandType = ApiCmd
				m.FillInfo()
				m.FindHashBlades()
			} else {
				if apiReply.Status[0].Status == "E" && strings.Contains(apiReply.Status[0].Msg, "Missing JSON") {
					apiCommand.SetCommand("summary", ApiCommand)
					rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand())
					json.Unmarshal(rawReply, &apiReply)
					if len(apiReply.Status) > 0 {
						m.ApiCommandType = ApiCommand
						switch apiReply.Status[0].Description {
						case "cgminer 4.9.2":
							m.Type = Whatsminer

						case "cgminer 1.0.0", "cgminer 4.9.0 rwglr", "jansson 2.12", "cgminer 4.9.0":
							m.Type = Antminer
						case "cgminer 4.11.1":
							apiCommand.SetCommand("stats", ApiCommand)
							rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand())
							if strings.Contains(string(rawReply), "Antminer") {
								m.Type = Antminer

							} else {
								m.Type = Avalon
							}
						case "cgminer 4.10.0":
							m.Type = Innosilicon
						default:
							m.Type = UnknownMiner
						}
						if m.Type != UnknownMiner {
							m.FillInfo()
							m.FindHashBlades()
						}

					}
				}
			}
		}
	}
}

func (m *miner) PrintInfo() {
	//u.NewUnit("Hashes per second", "H/s")
	fmt.Printf("IP Address: %s\n", m.IpAddress)
	fmt.Printf("MAC Address: %s\n", m.MacAddress)
	fmt.Printf("Model: %s\n", m.Model)
	fmt.Printf("Hashrate: %f\n", m.Hashrate)
	fmt.Printf("Status: %s\n", MinerStatusDesc[m.Status])
	fmt.Printf("Firmware: %s\n", m.FwVersion)
}

func (m *miner) FillInfo() {
	var (
		apiCommand MinerApiCommand
		rawReply   []byte
		err        error
	)
	apiCommand.SetCommand("summary", m.ApiCommandType)
	if rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand()); err != nil {
		switch cause := err.(type) {
		case net.Error:
			if cause.Timeout() {
				m.Status = Timeout
			}
		default:
			log.Println(err.Error())
		}
	} else {
		switch m.Type {
		case Whatsminer:
			replyWmSummary := new(WmReplySummary)
			replyWmPools := new(WmReplyPools)
			json.Unmarshal(rawReply, &replyWmSummary)
			m.MacAddress = replyWmSummary.Summary[0].MacAddress
			m.FwVersion = replyWmSummary.Summary[0].FwVersion
			m.Hashrate = replyWmSummary.Summary[0].HrAvg * 1000000
			m.Power = replyWmSummary.Summary[0].PowerRT
			m.Fans = append(m.Fans, minerFan{Id: "Fan Speed In", Speed: replyWmSummary.Summary[0].FanSpeedIn})
			m.Fans = append(m.Fans, minerFan{Id: "Fan Speed Out", Speed: replyWmSummary.Summary[0].FanSpeedOut})
			m.Fans = append(m.Fans, minerFan{Id: "Power Fanspeed", Speed: replyWmSummary.Summary[0].PowerFanspeed})
			apiCommand.SetCommand("devdetails", m.ApiCommandType)
			rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand())
			replyWmDevDetails := new(WmReplyDevDetails)
			json.Unmarshal(rawReply, &replyWmDevDetails)
			if len(replyWmDevDetails.Devs) > 0 {
				m.Model = replyWmDevDetails.Devs[0].Model
			}
			apiCommand.SetCommand("pools", m.ApiCommandType)
			rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand())
			log.Println(string(rawReply[:]))
			json.Unmarshal(rawReply, &replyWmPools)
			m.Pool1 = replyWmPools.Pools[0].URL
			m.Worker1 = replyWmPools.Pools[0].User
		case Antminer:
			replyAmSummary := new(AmReplySummary)
			json.Unmarshal(rawReply, &replyAmSummary)
			if len(replyAmSummary.Summary) > 0 {
				m.Hashrate = replyAmSummary.Summary[0].GHSAv * 1000000000
			} else {
				err := fmt.Errorf("IP: %s, empty summary", m.IpAddress)
				log.Println(err.Error())
			}
			replyAmStats := new(AmReplyStats)
			apiCommand.SetCommand("stats", m.ApiCommandType)
			if rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand()); err != nil {
				log.Println(err.Error())
			} else {
				rawReply = []byte(strings.Replace(string(rawReply), "}{", "},{", -1))
				json.Unmarshal(rawReply, &replyAmStats)
				if len(replyAmStats.Stats) > 0 {
					m.Model = replyAmStats.Stats[0].Type
					m.FwVersion = replyAmStats.Stats[0].Miner
					m.NumBlades = replyAmStats.Stats[1].MinerCount
					for i := 1; i <= replyAmStats.Stats[1].FanNum; i++ {
						fanId := fmt.Sprintf("Fan%d", i)
						r := reflect.ValueOf(&replyAmStats.Stats[1]).Elem()
						fanSpeed := r.FieldByName(fanId).Int()
						m.Fans = append(m.Fans, minerFan{
							Id:    fanId,
							Speed: int(fanSpeed),
						})
					}
				} else {
					err := fmt.Errorf("IP: %s, empty stats", m.IpAddress)
					log.Println(err.Error())
					log.Println(string(rawReply))
				}
			}
			replyAmPools := new(AmReplyPools)
			apiCommand.SetCommand("pools", m.ApiCommandType)
			rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand())
			json.Unmarshal(rawReply, &replyAmPools)
			m.Pool1 = replyAmPools.Pools[0].URL
			m.Worker1 = replyAmPools.Pools[0].User
		case Innosilicon:
			m.Model = "Innosilicon"

		case Avalon:
			m.Model = "Avalon"
		}
	}
}

func (m *miner) FindHashBlades() {
	var (
		apiCommand MinerApiCommand
		rawReply   []byte
		err        error
	)
	switch m.Type {
	case Whatsminer:
		apiCommand.SetCommand("devs", m.ApiCommandType)
		if rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand()); err != nil {

		} else {
			m.Status = MinerOK
			replyWmDevs := new(WmReplyDevs)
			json.Unmarshal(rawReply, &replyWmDevs)
			m.NumBlades = len(replyWmDevs.Devs)
			for i := 0; i < m.NumBlades; i++ {
				var alive bool
				if replyWmDevs.Devs[i].Status == "Alive" {
					alive = true
				} else {
					m.Status = HaveDeadHashBlades
				}
				m.HashBlades = append(m.HashBlades, minerHashblade{
					Alive:       alive,
					Hashrate:    int64(replyWmDevs.Devs[i].HrAvg),
					Temperature: int(replyWmDevs.Devs[i].Temperature),
				})
				if replyWmDevs.Devs[i].Slot != replyWmDevs.Devs[i].Id {
					m.Status = MissingHashBlade
				}
			}
		}
	case Antminer:
		//fmt.Println(AmGetDevs(m.IpAddress, m.ApiPort))
	}
}

func (m *miner) Test() {
	var (
		apiCommand MinerApiCommand
		rawReply   []byte
		err        error
	)
	log.Println("SUMMARY")
	apiCommand.SetCommand("summary", m.ApiCommandType)
	if rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand()); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(string(rawReply))
	}
	log.Println("STATS")
	apiCommand.SetCommand("stats", m.ApiCommandType)
	if rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand()); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(string(rawReply))
	}
	log.Println("DEVS")
	apiCommand.SetCommand("devs", m.ApiCommandType)
	if rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand()); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(string(rawReply))
	}
	log.Println("DEVDETAILS")
	apiCommand.SetCommand("devdetails", m.ApiCommandType)
	if rawReply, err = SendApiCommand(m.IpAddress, m.ApiPort, apiCommand.GetCommand()); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(string(rawReply))
	}

}

func SendApiCommand(ipaddress string, port int, command []byte) (res []byte, err error) {
	var (
		conn  net.Conn
		reply []byte
	)
	addr := fmt.Sprintf("%s:%d", ipaddress, port)
	if conn, err = net.DialTimeout("tcp", addr, time.Duration(3*time.Second)); err != nil {
		fmt.Println(err.Error())
	} else {
		defer conn.Close()
		if _, err = conn.Write(command); err != nil {
			fmt.Println(err.Error())
		} else {
			conn.(*net.TCPConn).CloseWrite()
			if reply, err = ioutil.ReadAll(conn); err != nil {
				fmt.Println(err.Error())
			} else {
				res = bytes.Trim(reply, "\x00")
			}
		}
	}
	return
}
