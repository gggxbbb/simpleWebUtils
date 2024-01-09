package main

import (
	"github.com/gin-gonic/gin"
	"simpleWebUtils/router/code"
	"simpleWebUtils/router/ip"
	"simpleWebUtils/router/minecraft"
	"simpleWebUtils/router/ua"
	"strconv"
)

var variablePorts string = "0"

func main() {

	r := gin.Default()

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
