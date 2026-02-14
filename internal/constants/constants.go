package constants

import (
	"errors"
	"fmt"
)

const (
	AssetDir = "./data"
)

var (
	ErrInvalidIP = errors.New("invalid IP")
)

const (
	MaxMindDownloadPathQuery   = "/app/geoip_download?edition_id=%s&license_key=%s&suffix=%s"
	DBArchiveDownloadSuffix    = "tar.gz"
	DBSHA256FileDownloadSuffix = "tar.gz.sha256"
	DBTypeCountry              = "GeoLite2-Country"
	DBTypeCity                 = "GeoLite2-City"
	DBTypeASN                  = "GeoLite2-ASN"
	DBSuffix                   = "mmdb"
)

var (
	MaxMindHost        = "https://download.maxmind.com"
	MaxMindDownloadURL = fmt.Sprintf("%s%s", MaxMindHost, MaxMindDownloadPathQuery)
)
