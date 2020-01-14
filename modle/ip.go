package modle

import (
	"apiTools/utils"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"path/filepath"
)

// ipv4表单结构
type Ipv4Form struct {
	Ip string `form:"ip" json:"ip" xml:"ip" binding:"required"`
}

// ipv4数据返回结构
type Ipv4Info struct {
	CityId   int64  `json:"cityId"`   // 城市ID号
	Country  string `json:"country"`  // 国家名称
	Region   string `json:"region"`   // 区域号
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	ISP      string `json:"isp"`      // isp厂商
}

var (
	ipv4db *ip2region.Ip2Region
)

func init() {
	// 初始化ipv4db数据库信息
	region, err := ip2region.New(filepath.Join(utils.GetRootPath(), "data/ip/ip2region.db"))
	if err != nil {
		panic(err)
	}
	ipv4db = region

}

// ipv4信息查询
func Ipv4Query(ipv4Form Ipv4Form) (Ipv4Info, error) {
	ipInfo := Ipv4Info{}
	queryIpInfo, err := ipv4db.MemorySearch(ipv4Form.Ip)
	if err != nil {
		return ipInfo, err
	}
	ipInfo.CityId = queryIpInfo.CityId

	if queryIpInfo.Country == "0" {
		queryIpInfo.Country = ""
	}
	ipInfo.Country = queryIpInfo.Country

	if queryIpInfo.Region == "0" {
		queryIpInfo.Region = ""
	}
	ipInfo.Region = queryIpInfo.Region

	if queryIpInfo.Province == "0" {
		queryIpInfo.Province = ""
	}
	ipInfo.Province = queryIpInfo.Province

	if queryIpInfo.City == "0" {
		queryIpInfo.City = ""
	}
	ipInfo.City = queryIpInfo.City

	if queryIpInfo.ISP == "0" {
		queryIpInfo.ISP = ""
	}
	ipInfo.ISP = queryIpInfo.ISP

	return ipInfo, nil
}
