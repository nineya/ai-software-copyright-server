package redbook

type RouterGroup struct {
	CookieApiRouter
	ProhibitedApiRouter
	RedbookApiRouter
	WriteApiRouter
}
