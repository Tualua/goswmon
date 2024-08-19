package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type jsonKeaRequest struct {
	Command string   `json:"command"`
	Service []string `json:"service"`
}

func (j *jsonKeaRequest) ToBytes() (res []byte) {
	res, _ = json.Marshal(j)
	return
}

func (j *jsonKeaRequest) ToString() (res string) {
	data, _ := json.Marshal(j)
	res = string(data)
	return
}

type jsonKeaResponse struct {
	Arguments struct {
		Leases []KeaDhcpLease `json:"leases"`
	} `json:"arguments"`
	Result int    `json:"result"`
	Text   string `json:"text"`
}

type KeaDhcpLease struct {
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
}

func invokeKeaCommand(api string, port int, command string, service []string) (res jsonKeaResponse, err error) {
	var (
		req          *http.Request
		response     *http.Response
		responseData []byte
		jsonData     []jsonKeaResponse
	)

	apiUrl := fmt.Sprintf("http://%s:%d", api, port)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := jsonKeaRequest{
		Command: command,
		Service: service,
	}

	if req, err = http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(body.ToBytes())); err != nil {
		log.Println(err.Error())
	} else {
		req.Header.Add("Content-Type", "application/json")

	}

	if response, err = client.Do(req); err != nil {
		log.Println(err.Error())
	} else {
		if responseData, err = ioutil.ReadAll(response.Body); err != nil {
			fmt.Println(err.Error())
		} else {
			json.Unmarshal(responseData, &jsonData)
			res = jsonData[0]
		}
	}

	return res, err
}

func KeaGetLeasesDhcp4(site Site) (res []KeaDhcpLease, err error) {
	var (
		data jsonKeaResponse
	)
	if data, err = invokeKeaCommand(site.DhcpServer, site.DhcpApiPort, "lease4-get-all", []string{"dhcp4"}); err != nil {
		log.Println(err.Error())
	} else {
		res = data.Arguments.Leases
	}
	return
}
