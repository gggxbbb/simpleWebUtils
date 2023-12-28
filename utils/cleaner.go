package utils

import (
	"regexp"
)

// LocalAddressCleaner replace any 127.0.0.1:<port>/172.x.x.x:<port>/10.x.x.x:<port> with <local-address>
func LocalAddressCleaner(ipt string) string {
	re := regexp.MustCompile(`(127\.0\.0\.1|172\.\d{1,3}\.\d{1,3}\.\d{1,3}|10\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d+)`)
	return re.ReplaceAllString(ipt, "<local-address>")
}
