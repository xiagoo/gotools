package ipx

import (
	"encoding/json"
	"fmt"
	"github.com/xiagoo/gotools/constx"
	"github.com/xiagoo/gotools/httpx"
	"math/rand"
	"strconv"
	"strings"
)

/*
https://www.opengps.cn/Data/IP/LocHighAcc.aspx
https://ip.rtbasia.com/
http://www.ipplus360.com/ip/
https://www.ipip.net/ip.html/
https://mall.ipplus360.com/
*/

const (
	AccuracyDistrict = "district"
	AccuracyLocate   = "locate"
)

type IpInfo struct {
	Msg      string `json:"msg"`
	Area     string `json:"area"`
	Code     string `json:"code"`
	Charge   bool   `json:"charge"`
	Data     *Data  `json:"data"`
	Ip       string `json:"ip"`
	Coordsys string `json:"coordsys"`
}

type Data struct {
	Continent string       `json:"continent"` //哪个洲
	Zipcode   string       `json:"zipcode"`
	Country   string       `json:"country"`
	Isp       string       `json:"isp"`
	Accuracy  string       `json:"accuracy"` //精确范围
	MultiArea []*MultiArea `json:"multiAreas"`
}
type MultiArea struct {
	Address  string `json:"address"`
	Lng      string `json:"lng"`
	Lat      string `json:"lat"`
	Radius   string `json:"radius"` //覆盖半径KM
	Prov     string `json:"prov"`
	City     string `json:"city"`
	District string `json:"district"`
}

func getAgent() *httpx.Agent {
	agent := httpx.NewAgent()
	agent.AddHeader(map[string]string{
		"User-Agent": constx.UserAgentList[rand.Intn(len(constx.UserAgentList))],
	})
	return agent
}

func GetIpInfo(ip, accuracyType string) (*IpInfo, error) {
	agent := httpx.NewAgent()
	agent.AddHeader(map[string]string{
		"User-Agent": constx.UserAgentList[rand.Intn(len(constx.UserAgentList))],
	})
	resp, err := agent.Get(fmt.Sprintf("https://mall.ipplus360.com/center/ip/api?ip=%s&type=%s&coordsys=BD09", ip, accuracyType)).GetResponseBody(nil)
	if err != nil {
		return nil, err
	}
	ipInfo := &IpInfo{}
	err = json.Unmarshal(resp, ipInfo)
	if err != nil {
		return nil, err
	}
	return ipInfo, nil
}

//https://www.123cha.com
func GetIpInfoBy123cha(ip string) ([]string, error) {
	ipInfos := []string{}
	for i := 0; i < 4; i++ {
		resp, err := getAgent().Post("https://www.123cha.com/ip/get.php").AddData(map[string]string{
			"ip":  ip,
			"job": strconv.Itoa(i),
		}).GetResponseBody(nil)
		if err != nil {
			continue
		}
		ipInfos = append(ipInfos, string(resp))
	}
	return ipInfos, nil
}

type SBInfo struct {
	Organization    string `json:"organization"`
	Longitude       string `json:"longitude"`
	City            string `json:"city"`
	TimeZone        string `json:"timezone"`
	Isp             string `json:"isp"`
	Offset          int    `json:"offset"	`
	Region          string `json:"region"`
	Asn             int    `json:"asn"`
	AsnOrganization string `json:"asn_organization"`
	Country         string `json:"country"`
	Ip              string `json:"ip"`
	Latitude        string `json:"latitude"`
	ContinentCode   string `json:"continent_code"`
	CountryCode     string `json:"country_code"`
	RegionCode      string `json:"region_code"`
}

//https://ip.sb/api/
func GetIpInfoBySB(ip string) (*SBInfo, error) {
	resp, err := getAgent().Get(fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)).GetResponseBody(nil)
	if err != nil {
		return nil, err
	}
	sbInfo := &SBInfo{}
	err = json.Unmarshal(resp, sbInfo)
	if err != nil {
		return nil, err
	}
	return sbInfo, nil
}
/*
获取本机ip
curl ip.sb
'http://ip.6655.com/ip.aspx',
'http://members.3322.org/dyndns/getip',
'http://icanhazip.com/',
'http://ident.me/',
'http://ipecho.net/plain',
'http://whatismyip.akamai.com/',
]
 */

func GetAddress(ip string) string {
	address := ""
	ipInfo, err := GetIpInfo(ip, AccuracyLocate)
	if err != nil {
		return address
	}
	address = ipInfo.Data.MultiArea[0].Address
	if strings.Contains(address, "*") {
		ipInfo, err = GetIpInfo(ip, AccuracyDistrict)
		if err != nil {
			return address
		}
		multiArea := ipInfo.Data.MultiArea[0]
		address = fmt.Sprintf("%s%s%s", multiArea.Prov, multiArea.City, multiArea.District)
	}

	if strings.Contains(address, "*") {
		ipInfo, err := GetIpInfoBy123cha(ip)
		if err != nil {
			return address
		}
		address = ipInfo[0]
	}
	return address
}
