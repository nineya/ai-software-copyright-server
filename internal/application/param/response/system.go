package response

type EnvironmentResponse struct {
	StartTime int64  `json:"startTime"`
	Database  string `json:"database"`
	Version   string `json:"version"`
	Mode      string `json:"mode"`
}

type SystemInfoResponse struct {
	OsName      string `json:"osName"`
	OsVersion   string `json:"osVersion"`
	OsArch      string `json:"osArch"`
	CpuNum      int    `json:"cpuNum"`
	TotalMemory uint64 `json:"totalMemory"`
	DiskSize    uint64 `json:"diskSize"`
	Timezone    string `json:"timezone"`
	UserName    string `json:"userName"`
	UserHome    string `json:"userHome"`
}

type UsageResponse struct {
	CpuUsageRate       float64 `json:"cpuUsageRate"`
	UsageMemory        uint64  `json:"usageMemory"`
	UsageDiskSize      uint64  `json:"usageDiskSize"`
	AllocPost          int     `json:"allocPost"`
	TotalAllocMemory   uint64  `json:"totalAllocMemory"`
	UsageAllocMemory   uint64  `json:"usageAllocMemory"`
	GCNum              uint32  `json:"gcNum"`
	UsageAllocDiskSize int64   `json:"usageAllocDiskSize"`
}

type AppResponse struct {
	Pid            int    `json:"pid"`
	StartTime      int64  `json:"startTime"`
	UpTime         int64  `json:"upTime"`
	WorkPath       string `json:"workPath"`
	StartPath      string `json:"startPath"`
	BuildGoVersion string `json:"buildGoVersion"`
	VersionCode    string `json:"versionCode"`
}

type SiteInfoResponse struct {
	SiteDiskSize int64  `json:"siteDiskSize"`
	SitePath     string `json:"sitePath"`
	SiteLogPath  string `json:"siteLogPath"`
}
