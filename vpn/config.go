package vpn

import "net/url"

func (vpn *VpnApi) DownloadConf(login string) []byte {
	answer, err := vpn.ApiFormDataRequest("api/user/config/show", url.Values{"username": []string{login}})
	if err != nil {
		return nil
	}
	return answer
}
