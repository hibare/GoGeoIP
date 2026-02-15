package maxmind

// DBType represents the type of MaxMind database.
type DBType string

const (
	DBTypeCountry DBType = "GeoLite2-Country"
	DBTypeCity    DBType = "GeoLite2-City"
	DBTypeASN     DBType = "GeoLite2-ASN"
)

// IPCountry represents country information for an IP.
type IPCountry struct {
	IP                  string `json:"ip"`
	Country             string `json:"country"`
	Continent           string `json:"continent"`
	ISOCountryCode      string `json:"iso_country_code"`
	ISOContinentCode    string `json:"iso_continent_code"`
	IsAnonymousProxy    bool   `json:"is_anonymous_proxy"`
	IsSatelliteProvider bool   `json:"is_satellite_provider"`
}

// IPCity represents city information for an IP.
type IPCity struct {
	City string `json:"city"`
	IPCountry
	Timezone  string  `json:"timezone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// IPASN represents ASN information for an IP.
type IPASN struct {
	IP           string `json:"ip"`
	ASN          uint   `json:"asn"`
	Organization string `json:"organization"`
}

// GeoIP represents complete geographic and ASN information for an IP.
type GeoIP struct {
	IPCity
	IPASN
	IP     string `json:"ip"`
	Remark string `json:"remark,omitempty"`
}
