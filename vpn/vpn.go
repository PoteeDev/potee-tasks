package vpn

import (
	"encoding/base64"
	"fmt"
	"os"
)

type VpnApi struct {
	Url       string
	BasicAuth string
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func CreateVPNClient() *VpnApi {
	vpn := VpnApi{
		Url:       fmt.Sprintf("http://%s", os.Getenv("OVPN_HOST")),
		BasicAuth: basicAuth(os.Getenv("OVPN_USER"), os.Getenv("OVPN_PASSWORD")),
	}
	return &vpn
}
