package vpn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (vpn *VpnApi) ApiFormDataRequest(route string, form url.Values) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", vpn.Url, route)
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+vpn.BasicAuth)
	if err != nil {
		return nil, err
	}
	clinet := http.Client{}
	resp, err := clinet.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (vpn *VpnApi) ApiJsonRequest(route string, data interface{}) ([]byte, error) {
	clentData := new(bytes.Buffer)
	json.NewEncoder(clentData).Encode(data)
	url := fmt.Sprintf("%s/%s", vpn.Url, route)
	req, err := http.NewRequest("POST", url, clentData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+vpn.BasicAuth)
	if err != nil {
		return nil, err
	}
	clinet := http.Client{}
	resp, err := clinet.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (vpn *VpnApi) CreateUser(login, password string) (string, error) {
	answer, err := vpn.ApiFormDataRequest("api/user/create", url.Values{
		"username": []string{login},
		"password": []string{password},
	})
	return string(answer), err
}

type CcdClient struct {
	User          string
	ClientAddress string
	CustomRoutes  []CustomRoute
}
type CustomRoute struct {
	Address string
	Mask    string
}

func CidrToRoute(cidr string) CustomRoute {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal(err)
	}
	return CustomRoute{
		ip.String(),
		net.IP(network.Mask).String(),
	}
}

func (vpn *VpnApi) AddUserRoutes(login, cidr string) (string, error) {
	answer, err := vpn.ApiJsonRequest("api/user/ccd/apply", CcdClient{
		login,
		"dynamic",
		[]CustomRoute{CidrToRoute(cidr)},
	})
	return string(answer), err
}

func (vpn *VpnApi) AddClient(login, password, cidr string) error {
	if _, err := vpn.CreateUser(login, password); err != nil {
		return err
	}
	if _, err := vpn.AddUserRoutes(login, cidr); err != nil {
		return err
	}
	return nil
}

func (vpn *VpnApi) UpdateClient() error {
	return nil
}

func (vpn *VpnApi) DeleteClient() error {
	return nil
}
