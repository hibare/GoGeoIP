package maxmind

import (
	"net"

	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/oschwald/geoip2-golang"
)

// Country looks up country information for an IP.
func (c *Client) Country(ip net.IP) (*geoip2.Country, error) {
	c.mu.RLock()
	reader, ok := c.readers[DBTypeCountry]
	c.mu.RUnlock()

	if !ok || reader == nil {
		return nil, ErrCountryDBNotLoaded
	}
	return reader.Country(ip)
}

// City looks up city information for an IP.
func (c *Client) City(ip net.IP) (*geoip2.City, error) {
	c.mu.RLock()
	reader, ok := c.readers[DBTypeCity]
	c.mu.RUnlock()

	if !ok || reader == nil {
		return nil, ErrCityDBNotLoaded
	}
	return reader.City(ip)
}

// ASN looks up ASN information for an IP.
func (c *Client) ASN(ip net.IP) (*geoip2.ASN, error) {
	c.mu.RLock()
	reader, ok := c.readers[DBTypeASN]
	c.mu.RUnlock()

	if !ok || reader == nil {
		return nil, ErrASNDBNotLoaded
	}
	return reader.ASN(ip)
}

// IP2Country looks up country information for an IP address.
func (c *Client) IP2Country(ipStr string) (IPCountry, error) {
	ipCountry := IPCountry{}

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return ipCountry, constants.ErrInvalidIP
	}

	record, err := c.Country(parsedIP)
	if err != nil {
		return ipCountry, err
	}

	ipCountry.IP = ipStr
	ipCountry.Continent = record.Continent.Names["en"]
	ipCountry.Country = record.Country.Names["en"]
	ipCountry.ISOContinentCode = record.Continent.Code
	ipCountry.ISOCountryCode = record.Country.IsoCode
	ipCountry.IsAnonymousProxy = record.Traits.IsAnonymousProxy
	ipCountry.IsSatelliteProvider = record.Traits.IsSatelliteProvider

	return ipCountry, nil
}

// IP2City looks up city information for an IP address.
func (c *Client) IP2City(ipStr string) (IPCity, error) {
	ipCity := IPCity{}

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return ipCity, constants.ErrInvalidIP
	}

	record, err := c.City(parsedIP)
	if err != nil {
		return ipCity, err
	}

	ipCity.IP = ipStr
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

// IP2ASN looks up ASN information for an IP address.
func (c *Client) IP2ASN(ipStr string) (IPASN, error) {
	ipAsn := IPASN{}

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return ipAsn, constants.ErrInvalidIP
	}

	record, err := c.ASN(parsedIP)
	if err != nil {
		return ipAsn, err
	}

	ipAsn.IP = ipStr
	ipAsn.ASN = record.AutonomousSystemNumber
	ipAsn.Organization = record.AutonomousSystemOrganization

	return ipAsn, nil
}

// IP2Geo looks up all geographic information for an IP address.
func (c *Client) IP2Geo(ipStr string) (GeoIP, error) {
	geoIP := GeoIP{
		IP: ipStr,
	}

	ipCity, err := c.IP2City(ipStr)
	if err != nil {
		return geoIP, err
	}

	ipAsn, err := c.IP2ASN(ipStr)
	if err != nil {
		return geoIP, err
	}

	geoIP.IPCity = ipCity
	geoIP.IPASN = ipAsn

	return geoIP, nil
}
