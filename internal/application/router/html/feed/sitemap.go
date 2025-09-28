package feed

import (
	_xml "ai-software-copyright-server/internal/application/model/xml"
	"ai-software-copyright-server/internal/application/param/response"
	slSev "ai-software-copyright-server/internal/application/service/short_link"
	"ai-software-copyright-server/internal/global"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SitemapXml(c *gin.Context) {
	sitemapData := getSitemapData()
	list := make([]_xml.SitemapUrl, 0)
	// 首页
	createSiteTime := time.Date(2019, 11, 28, 20, 5, 2, 0, time.UTC)
	list = append(list, _xml.SitemapUrl{global.Host, &createSiteTime})
	// 短链
	for _, netdisk := range sitemapData.Netdisk {
		list = append(list, _xml.SitemapUrl{global.Host + "/s/" + netdisk.Alias, netdisk.CreateTime})
	}

	c.Writer.Write([]byte(xml.Header))
	c.XML(http.StatusOK, _xml.SitemapElement{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Url:   list,
	})
}

func SitemapHtml(c *gin.Context) {
	htmlResponse := response.GenerateHtmlResult(c)
	htmlResponse.OkWithData("internal/feed/sitemap.html", getSitemapData())
}

func getSitemapData() response.SitemapResponse {
	netdisk, _ := slSev.GetNetdiskService().GetByLast(2000)

	return response.SitemapResponse{
		Netdisk: netdisk,
	}
}
