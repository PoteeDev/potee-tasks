package handlers

import (
	"console/vpn"
	"net/http"
)

var vpnApi = vpn.CreateVPNClient()

func (h *handler) DownloadVpn(w http.ResponseWriter, r *http.Request) {
	name, _ := h.GetUserName(r)
	configFile := vpnApi.DownloadConf(name)
	if configFile != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(configFile)
}
