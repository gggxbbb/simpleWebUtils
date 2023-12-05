package main

import (
	"regexp"
)

// replace any 127.0.0.1:<port> with <local-address>
func localAddressCleaner(ipt string) string {
	re := regexp.MustCompile(`127\.0\.0\.1:\d+`)
	return re.ReplaceAllString(ipt, "<local-address>")
}
