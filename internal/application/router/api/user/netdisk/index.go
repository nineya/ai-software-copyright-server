package netdisk

type RouterGroup struct {
	HelperApiRouter
	NetdiskApiRouter
	ResourceApiRouter
	SearchAppApiRouter
	SearchSiteApiRouter
	SearchWxampApiRouter
	ShortLinkApiRouter
}
