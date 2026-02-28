package maxmind

import "fmt"

const (
	MaxMindDownloadPathQuery   = "/app/geoip_download?edition_id=%s&license_key=%s&suffix=%s"
	DBArchiveDownloadSuffix    = "tar.gz"
	DBSHA256FileDownloadSuffix = "tar.gz.sha256"
	DBSuffix                   = "mmdb"
)

var (
	MaxMindHost        = "https://download.maxmind.com"
	MaxMindDownloadURL = fmt.Sprintf("%s%s", MaxMindHost, MaxMindDownloadPathQuery)
)
