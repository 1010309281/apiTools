package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modle"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Api文档首页
func ApiIndex(c *gin.Context) {
	allApiInfo, err := modle.GetAllApiInfo()
	if err != nil {
		logger.Echo.WithFields(logrus.Fields{
			"router": c.Request.URL.Path,
			"err":    err,
			"info":   allApiInfo,
		}).Error("access api index page fild")
		c.HTML(http.StatusBadGateway, "error.html", gin.H{"errorMsg": "系统异常请稍后重试"})
		return
	}
	c.HTML(http.StatusOK, "index.html", allApiInfo)
}

// Api文档详情页
func ApiDocs(c *gin.Context) {
	var apiDocFileName, urlPath, countKey string

	// 获取数据
	reqApiName := c.Param("apiName")
	urlPath = c.GetString("urlPath")
	if reqApiName != urlPath {
		urlPath = ""
	}
	apiDocFileName = c.GetString("docFile")
	countKey = c.GetString("countKey")

	if apiDocFileName == "" || urlPath == "" || countKey == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}
	apiDocInfo, err := modle.GetApiDocInfo(apiDocFileName, urlPath, countKey)
	if err != nil {
		logger.Echo.WithFields(logrus.Fields{
			"router": c.Request.URL.Path,
			"err":    err,
			"info":   apiDocInfo,
		}).Error("access api doc page fild")
		c.HTML(http.StatusBadGateway, "error.html", gin.H{"errorMsg": "系统异常请稍后重试"})
		return
	}
	c.HTML(http.StatusOK, "docs.html", apiDocInfo)
}


// 关于ApiTools
func ApiAbout(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}