package ip

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

type data struct {
	IP      string
	OK      bool
	Country string
	City    string
	ISOCode string
	Detail  detail
}

type dCity struct {
	Name      string
	GeoNameID uint
}

type dPostal struct {
	Code string
}

type dContinent struct {
	Name      string
	Code      string
	GeoNameID uint
}

type dCountry struct {
	Name      string
	ISOCode   string
	GeoNameID uint
	IsInEU    bool
}

type dTraits struct {
	IsAnonymousProxy    bool
	IsSatelliteProvider bool
}

type detail struct {
	City      dCity
	Postal    dPostal
	Continent dContinent
	Country   dCountry
	Traits    dTraits
}

func analyzeIP(ip string) data {

	db, err := geoip2.Open("./GeoIP/GeoLite2-City.mmdb")
	if err != nil {
		return data{IP: ip, OK: false}
	}
	defer db.Close()

	record, err := db.City(net.ParseIP(ip))

	if err != nil {
		return data{IP: ip, OK: false}
	}

	return data{
		IP:      ip,
		OK:      true,
		Country: record.Country.Names["en"],
		City:    record.City.Names["en"],
		ISOCode: record.Country.IsoCode,
		Detail: detail{
			City: dCity{
				Name:      record.City.Names["en"],
				GeoNameID: record.City.GeoNameID,
			},
			Postal: dPostal{
				Code: record.Postal.Code,
			},
			Continent: dContinent{
				Name:      record.Continent.Names["en"],
				Code:      record.Continent.Code,
				GeoNameID: record.Continent.GeoNameID,
			},
			Country: dCountry{
				Name:      record.Country.Names["en"],
				ISOCode:   record.Country.IsoCode,
				GeoNameID: record.Country.GeoNameID,
				IsInEU:    record.Country.IsInEuropeanUnion,
			},
			Traits: dTraits{
				IsAnonymousProxy:    record.Traits.IsAnonymousProxy,
				IsSatelliteProvider: record.Traits.IsSatelliteProvider,
			},
		},
	}

}
