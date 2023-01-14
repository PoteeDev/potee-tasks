package config

import (
	"console/models"
	"fmt"
	"net"

	"github.com/lib/pq"
)

func GenerateCidr(firsIp, LastIp string) string {
	ip1 := net.ParseIP(firsIp)
	ip2 := net.ParseIP(LastIp)
	maxLen := 32
	for l := maxLen; l >= 0; l-- {
		mask := net.CIDRMask(l, maxLen)
		na := ip1.Mask(mask)
		n := net.IPNet{IP: na, Mask: mask}

		if n.Contains(ip2) {
			// fmt.Printf("smallest possible CIDR range: %v/%v\n", na, l)
			return fmt.Sprintf("%v/%v", na, l)
		}
	}
	return ""
}

func (c *Config) GenerateClientPools() *[]models.Pool {
	allHosts, _ := Hosts(c.ServicesSubnet)
	//remove network address from pool
	var newHosts []string
	for _, ip := range allHosts {
		lastOctet := int(net.ParseIP(ip).To4()[3])
		if lastOctet%(c.IpCount+2) == 0 {
			continue
		}
		newHosts = append(newHosts, ip)
	}

	var ipPools = []models.Pool{}
	ipList := pq.StringArray{}
	for i, ip := range newHosts {
		ipList = append(ipList, ip)
		if (i+1)%(c.IpCount+1) == 0 && i > 0 {
			pool := models.Pool{
				IPPool: ipList,
			}
			pool.Cidr = GenerateCidr(ipList[0], ipList[len(ipList)-1])
			ipPools = append(ipPools, pool)
			ipList = pq.StringArray{}
		}
	}
	return &ipPools
}
func (c *Config) GenerateClientAvaliableIP() []*models.AvaliableIP {
	allocatedIPs, _ := Hosts(c.ClinetsSubnet)

	var ips []*models.AvaliableIP
	for _, ip := range allocatedIPs[2:] {
		ips = append(ips, &models.AvaliableIP{Ip: ip})
	}
	return ips

}
func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil

	default:
		return ips[1 : len(ips)-1], nil
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
