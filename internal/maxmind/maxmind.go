package maxmind

import (
	"net"

	"github.com/hibare/GoGeoIP/internal/constants"
)

type IPCountry struct {
	IP                  string `json:"ip"`
	Country             string `json:"country"`
	Continent           string `json:"continent"`
	ISOCountryCode      string `json:"iso_country_code"`
	ISOContinentCode    string `json:"iso_continent_code"`
	IsAnonymousProxy    bool   `json:"is_anonymous_proxy"`
	IsSatelliteProvider bool   `json:"is_satellite_provider"`
}

type IPCity struct {
	IP   string `json:"ip"`
	City string `json:"city"`
	IPCountry
	Timezone  string  `json:"timezone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type IPASN struct {
	IP           string `json:"ip"`
	ASN          uint   `json:"asn"`
	Organization string `json:"oraganization"`
}

type GeoIP struct {
	IPCity
	IPASN
	IP     string `json:"ip"`
	Remark string `json:"remark,omitempty"`
}

func IP2Country(ip string) (IPCountry, error) {
	ipCountry := IPCountry{}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ipCountry, constants.ErrInvalidIP
	}

	record, err := GetDB(constants.DBTypeCountry).Country(parsedIP)
	if err != nil {
		return ipCountry, err
	}

	ipCountry.IP = ip
	ipCountry.Continent = record.Continent.Names["en"]
	ipCountry.Country = record.Country.Names["en"]
	ipCountry.ISOContinentCode = record.Continent.Code
	ipCountry.ISOCountryCode = record.Country.IsoCode
	ipCountry.IsAnonymousProxy = record.Traits.IsAnonymousProxy
	ipCountry.IsSatelliteProvider = record.Traits.IsSatelliteProvider

	return ipCountry, nil
}

func IP2City(ip string) (IPCity, error) {
	ipCity := IPCity{}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ipCity, constants.ErrInvalidIP
	}

	record, err := GetDB(constants.DBTypeCity).City(parsedIP)
	if err != nil {
		return ipCity, err
	}

	ipCity.IP = ip
	ipCity.City = record.City.Names["en"]
	ipCity.Timezone = record.Location.TimeZone
	ipCity.Latitude = record.Location.Latitude
	ipCity.Longitude = record.Location.Longitude
	ipCity.Country = record.Country.Names["en"]
	ipCity.Continent = record.Continent.Names["en"]
	ipCity.ISOCountryCode = record.Country.IsoCode
	ipCity.ISOContinentCode = record.Continent.Code
	ipCity.IsAnonymousProxy = record.Traits.IsAnonymousProxy
	ipCity.IsSatelliteProvider = record.Traits.IsSatelliteProvider

	return ipCity, nil
}

func IP2ASN(ip string) (IPASN, error) {
	ipAsn := IPASN{}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ipAsn, constants.ErrInvalidIP
	}

	record, err := GetDB(constants.DBTypeASN).ASN(parsedIP)
	if err != nil {
		return ipAsn, err
	}

	ipAsn.IP = ip
	ipAsn.ASN = record.AutonomousSystemNumber
	ipAsn.Organization = record.AutonomousSystemOrganization

	return ipAsn, nil
}

func IP2Geo(ip string) (GeoIP, error) {
	geoIP := GeoIP{
		IP: ip,
	}

	ipCity, err := IP2City(ip)
	if err != nil {
		return geoIP, err
	}

	ipAsn, err := IP2ASN(ip)
	if err != nil {
		return geoIP, err
	}

	geoIP.IPCity = ipCity
	geoIP.IPASN = ipAsn

	return geoIP, nil
}
