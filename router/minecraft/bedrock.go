package minecraft

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net"
	"simpleWebUtils/utils"
	"strconv"
	"strings"
	"time"
)

type minecraftBedrockRemote struct {
	Server string `json:"server"`
	Port   string `json:"port"`
}

// motd-bedrock
func utilsMinecraftBedrock(ctx *gin.Context) {

	var server minecraftBedrockRemote

	if ctx.Request.Method == "GET" {
		//GET
		server = minecraftBedrockRemote{
			Server: ctx.Param("server"),
			Port:   ctx.Param("port"),
		}
	} else if ctx.Request.Method == "POST" {
		//POST
		err := ctx.Bind(&server)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":       "cannot parse request",
				"description": utils.LocalAddressCleaner(err.Error()),
			})
			return
		}
	}

	if server.Port == "" {
		server.Port = "19132"
	}

	//reserve server
	addr, err := net.ResolveUDPAddr("udp", server.Server+":"+server.Port)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot resolve server",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}

	target_ip := addr.IP.String()

	//connect to server
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot connect to server",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	//send payload
	payload, _ := hex.DecodeString("0100000000240D12D300FFFF00FEFEFEFEFDFDFDFD12345678")
	_, err = conn.Write(payload)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot send payload",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}

	//receive response
	buf := make([]byte, 1024)
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot set read deadline",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}
	_, err = conn.Read(buf)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot receive response",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}

	motd := string(buf)
	data := strings.Split(motd, ";")
	//remove last empty string
	data = data[:len(data)-1]

	online, err := strconv.Atoi(data[4])
	max, err := strconv.Atoi(data[5])
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot parse response",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"raw":       data,
		"server":    server.Server,
		"server_ip": target_ip,
		"port":      server.Port,
		"motd":      data[1],
		"protocol":  data[2],
		"version":   data[3],
		"online":    online,
		"max":       max,
		"level":     data[6],
		"levelName": data[7],
		"gamemode":  data[8],
	})
}
