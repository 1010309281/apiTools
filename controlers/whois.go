package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modle"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 域名whois信息查询 --> api
func WhoisQuery(c *gin.Context) {
	var whoisForm modle.WhoisForm
	err := c.Bind(&whoisForm)
	data := gin.H{
		"data":   "", // whois数据
		"msg":    "", // 域名查询状态
		"status": 5,  // 域名查询状态
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": data})
		logger.Echo.WithFields(logrus.Fields{
			"router": c.Request.URL.Path,
			"err":    err,
			"query": whoisForm,
		}).Error("get query whois form param fail")
		return
	}
	whoisInfo := &modle.WhoisInfo{}
	if whoisForm.OutType == "text" {
		whoisInfo.WhoisForm.OutType = "text"
		whoisInfo, err = modle.QueryWhoisInfo(&whoisForm)
		data["status"] = whoisInfo.Status
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 0, "data": data})
			logger.Echo.WithFields(logrus.Fields{
				"router": c.Request.URL.Path,
				"err":    err,
				"query":  whoisForm,
				"data":   whoisInfo,
			}).Error("query whois info fild")
			return
		}
		data["data"] = whoisInfo.TextInfo
	} else {
		whoisInfo.WhoisForm.OutType = "json"
		whoisInfo, err = modle.QueryWhoisInfoToJson(&whoisForm)
		data["status"] = whoisInfo.Status
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 0, "data": data})
			logger.Echo.WithFields(logrus.Fields{
				"router": c.Request.URL.Path,
				"err":    err,
				"query":  whoisForm,
				"data":   whoisInfo,
			}).Error("query whois info fild")
			return
		}
		data["data"] = whoisInfo.JsonInfo
	}
	// log debug
	logger.Echo.WithFields(logrus.Fields{
		"router": c.Request.URL.Path,
		"warn":   err,
		"query":  whoisForm,
		"data":   whoisInfo,
	}).Debug("query whois info success")
	// log info
	logger.Echo.WithFields(logrus.Fields{
		"router": c.Request.URL.Path,
		"warn":   err,
		"query":  whoisForm,
		"data":   logrus.Fields{"domain": whoisInfo.Domain},
	}).Info("query whois info success")

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}
