package handlers

import (
	"console/vpn"
	"net/http"
)

var vpnApi = vpn.CreateVPNClient()

func (h *handler) DownloadVpn(w http.ResponseWriter, r *http.Request) {
	user, _ := h.tk.ExtractTokenMetadata(r)
	configFile := vpnApi.DownloadConf(user.UserId)
	if configFile != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(configFile)
}
