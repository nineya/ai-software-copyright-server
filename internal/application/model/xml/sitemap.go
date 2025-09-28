package _xml

import (
	"encoding/xml"
	"time"
)

type SitemapElement struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	Url     []SitemapUrl `xml:"url"`
}

type SitemapUrl struct {
	Loc     string     `xml:"loc"`
	Lastmod *time.Time `xml:"lastmod"`
}
