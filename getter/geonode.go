package getter

import (
	"encoding/json"
	"fmt"
	"sync"
	"vote/utils"

	"github.com/henson/proxypool/pkg/models"
)

type IpData struct {
	Id             string   `json:"_id"`
	Ip             string   `json:"ip"`
	AnonymityLevel string   `json:"anonymityLevel"`
	Asn            string   `json:"asn"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
	Created_at     string   `json:"created_at"`
	Google         string   `json:"google"`
	Isp            string   `json:"isp"`
	LastChecked    string   `json:"lastChecked"`
	Latency        string   `json:"latency"`
	Org            string   `json:"org"`
	Port           string   `json:"port"`
	Protocols      []string `json:"protocols"`
	Region         string   `json:"region"`
	ResponseTime   int      `json:"responseTime"`
	Speed          int      `json:"speed"`
	Updated_at     string   `json:"updated_at"`
}
type IpRequestParam struct {
	Data  []IpData `json:"data"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
}

var geonodePage int = 1

// feiyi get ip from feiyiproxy.com
func Geonode() (result []*models.IP) {
	var wg sync.WaitGroup
	url := fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?limit=500&page=%d&sort_by=lastChecked&sort_type=desc", geonodePage)
	sendRequest := utils.NewSendRequest(nil, "")
	body, _, _ := sendRequest.Get(url, nil)
	response := &IpRequestParam{}
	json.Unmarshal(body, &response)
	for _, v := range response.Data {
		wg.Add(1)
		proxyAddr := fmt.Sprintf("socks5://%s:%s", v.Ip, v.Port)
		go func(proxyAddr string) {
			defer wg.Done()
			IP := models.NewIP()
			IP.Data = v.Ip + ":" + v.Port
			IP.Type1 = "socks5"
			IP.Source = "proxylist.geonode.com"
			result = append(result, IP)
			geonodePage++
		}(proxyAddr)
	}
	wg.Wait()
	return
}
