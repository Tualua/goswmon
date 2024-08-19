package main

import (
	"encoding/json"
	"fmt"
)

type WmToken struct {
	IpAddress string
	Port      int
}

type WmCommand struct {
	Cmd     string `json:"cmd"`
	Command string `json:"command,omitempty"`
}

func (c *WmCommand) SetCommand(command string, old bool) {
	if old {
		c.Command = command
	} else {
		c.Cmd = command
	}
}

func (c *WmCommand) GetCommand() (res []byte) {
	var (
		err error
	)
	if res, err = json.Marshal(c); err != nil {
		fmt.Println(err.Error())
	}
	return
}

func (c *WmCommand) ToString(command string, old bool) (res string) {
	if jsonCommand, err := json.Marshal(c); err != nil {
		fmt.Println(err.Error())
	} else {
		res = string(jsonCommand)
	}
	return
}

type WmSummary struct {
	HrAvg         float64 `json:"MHS av"`
	HrTarget      float64 `json:"Target MHS"`
	HrFactory     int64   `json:"Factory GHS"`
	MacAddress    string  `json:"MAC"`
	ChipData      string  `json:"Chip Data"`
	FwVersion     string  `json:"Firmware Version"`
	PowerRT       int     `json:"Power_RT"`
	FanSpeedIn    int     `json:"Fan Speed In"`
	FanSpeedOut   int     `json:"Fan Speed Out"`
	PowerFanspeed int     `json:"Power Fanspeed"`
}

type WmPool struct {
	Pool                int    `json:"POOL"`
	URL                 string `json:"URL"`
	Status              string `json:"Status"`
	User                string `json:"User"`
	CurrentBlockVersion int    `json:"Current Block Version"`
}

type WmDevs struct {
	Status      string  `json:"Status"`
	Id          int     `json:"ID"`
	Slot        int     `json:"Slot"`
	HrAvg       float64 `json:"MHS av"`
	Temperature float64 `json:"Temperature"`
}

type WmDevDetails struct {
	Id    int    `json:"ID"`
	Model string `json:"Model"`
}

type WmReplySummary struct {
	Status  []CgMinerStatus `json:"STATUS"`
	Summary []WmSummary     `json:"SUMMARY"`
}

type WmReplyDevs struct {
	Status []CgMinerStatus `json:"STATUS"`
	Devs   []WmDevs        `json:"DEVS"`
}

type WmReplyDevDetails struct {
	Status []CgMinerStatus `json:"STATUS"`
	Devs   []WmDevDetails  `json:"DEVDETAILS"`
}

type WmReplyPools struct {
	Status []CgMinerStatus `json:"STATUS"`
	Pools  []WmPool        `json:"POOLS"`
}
