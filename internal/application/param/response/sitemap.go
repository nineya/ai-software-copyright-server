package response

import "ai-software-copyright-server/internal/application/model/table"

type SitemapResponse struct {
	Netdisk []table.NetdiskShortLink `json:"netdisk"`
}
