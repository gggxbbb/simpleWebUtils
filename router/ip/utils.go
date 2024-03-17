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
	}

}
