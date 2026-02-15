export interface GeoIP {
  ip: string
  city: string
  country: string
  continent: string
  iso_country_code: string
  iso_continent_code: string
  is_anonymous_proxy: boolean
  is_satellite_provider: boolean
  timezone: string
  latitude: number
  longitude: number
  asn: number
  organization: string
  remark?: string
}

export interface LookupHistory {
  id: string
  ip: string
  organization: string
  location: string
  timestamp: Date
}
