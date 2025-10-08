package global

import (
	"regexp"
	"runtime"
	"strconv"
)

const SiteName = "MIDDLEWARE_SITE_NAME"
const DefaultThemeId = "default"
const AuthToken = "Auth"
const RefreshToken = "Refresh"
const Admin = "Admin"
const User = "User"

var (
	GitCommitHash  = "f835ebee06eb230eb9673f6f2af88b909f88f109"
	Version        = "0.0.1"
	tBuildTime     = "1705399795000"
	BuildTime, _   = strconv.ParseInt(tBuildTime, 10, 64)
	BuildGoVersion = runtime.Version()
	Host           = "http://127.0.0.1:9300"
	BotReg, _      = regexp.Compile("(okhttp|Go-http-client|python-requests|spider|bot)")
	AesKey         = []byte("6c088f40caea4e4bb9bd5403d7134783")
)
