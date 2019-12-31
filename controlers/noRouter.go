package controlers

import (
	modles "apiTools/modle"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 404处理(404跳转, 短链接跳转)
func NoRouter(c *gin.Context) {
	// 处理短链接跳转
	if len(c.Request.URL.Query()) > 0 {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}
	host, path := c.Request.Host, c.Request.URL.Path
	shortUrl := fmt.Sprintf("http://%s%s", host, path)
	shortInfo, err := modles.ParseShort(shortUrl)
	if err != nil {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}
	c.Redirect(http.StatusFound, shortInfo.LongUrl)
}
