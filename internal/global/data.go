package global

import (
	"runtime"
	"strconv"
)

const SiteName = "MIDDLEWARE_SITE_NAME"
const DefaultThemeId = "default"
const AuthToken = "Auth"
const RefreshToken = "Refresh"

var (
	GitCommitHash  = "f835ebee06eb230eb9673f6f2af88b909f88f109"
	Version        = "0.0.1"
	tBuildTime     = "1705399795000"
	BuildTime, _   = strconv.ParseInt(tBuildTime, 10, 64)
	BuildGoVersion = runtime.Version()
)
