package main

import (
	"encoding/json"
	"net/http"
)

func (m *miner) WriteResponse(w *http.ResponseWriter) {
	enc := json.NewEncoder(*w)
	enc.SetIndent("", "    ")
	enc.Encode(m)
}

func GetMinerInfo(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.Header.Get("X-Miner-IPAddress")
	if ipAddress != "" {
		m := new(miner)
		m.Init(ipAddress)
		m.WriteResponse(&w)
	}
}
