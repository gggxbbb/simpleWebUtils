package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"simpleWebUtils/router/code"
	"simpleWebUtils/router/ip"
	"simpleWebUtils/router/minecraft"
	"simpleWebUtils/router/ua"
	"strconv"
	"strings"
)

// using `go build -ldflags "-X main.variablePorts=1"` to enable variable ports
var variablePorts = "0"

func main() {

	r := gin.Default()
	r.RemoteIPHeaders = []string{
		"CF-Connecting-IP",
		"True-Client-IP",
		"X-Client-IP",
		"X-Original-Forwarded-For",
		"X-Forwarded-For",
		"X-Real-IP",
	}
	trustedProxies := []string{"127.0.0.1", "::1"}
	if envTrustedProxies := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES")); envTrustedProxies != "" {
		configuredTrustedProxies := make([]string, 0)
		for _, trustedProxy := range strings.Split(envTrustedProxies, ",") {
			trustedProxy = strings.TrimSpace(trustedProxy)
			if trustedProxy == "" {
				continue
			}
			configuredTrustedProxies = append(configuredTrustedProxies, trustedProxy)
		}
		trustedProxies = configuredTrustedProxies
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		panic("failed to configure trusted proxies: " + err.Error())
	}

	r.GET("/", readme)

	ua.Init(r)
	ip.Init(r)
	code.Init(r)
	minecraft.Init(r)

	port := 4399
	err := r.Run(":" + strconv.Itoa(port))
	if err != nil {
		if //goland:noinspection GoBoolExpressions
		variablePorts == "1" {
			for i := 0; i < 100; i++ {
				port++
				err = r.Run(":" + strconv.Itoa(port))
				if err == nil {
					print("Server started at port " + strconv.Itoa(port) + "\n")
					break
				}
			}
		} else {
			panic(err)
		}
	} else {
		print("Server started at port " + strconv.Itoa(port) + "\n")
	}
}
