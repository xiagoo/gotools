package ipx

import (
	"encoding/json"
	"fmt"
	"github.com/xiagoo/gotools/constx"
	"github.com/xiagoo/gotools/httpx"
	"math/rand"
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

func GetIpInfo(ip, accuracyType string) (*IpInfo, error) {
	agent := httpx.NewAgent()
	agent.AddHeader(map[string]string{
		"User-Agent":constx.UserAgentList[rand.Intn(len(constx.UserAgentList))],
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
	return address
}
