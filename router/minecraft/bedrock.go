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
		err := ctx.ShouldBindJSON(&server)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":       "invalid JSON request body",
				"description": utils.LocalAddressCleaner(err.Error()),
			})
			return
		}
	}

	if strings.TrimSpace(server.Server) == "" {
		ctx.JSON(400, gin.H{
			"error": "server is required",
		})
		return
	}

	if server.Port == "" {
		server.Port = "19132"
	}

	portInt, err := strconv.Atoi(server.Port)
	if err != nil || portInt < 1 || portInt > 65535 {
		ctx.JSON(400, gin.H{
			"error": "port must be a number between 1 and 65535",
		})
		return
	}

	//reserve server
	target := net.JoinHostPort(strings.Trim(server.Server, "[]"), server.Port)
	addr, err := net.ResolveUDPAddr("udp", target)
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
	n, err := conn.Read(buf)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot receive response",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}

	motd := string(buf[:n])
	if idx := strings.Index(motd, "MCPE;"); idx >= 0 {
		motd = motd[idx:]
	}
	motd = strings.TrimRight(motd, "\x00")
	data := strings.Split(motd, ";")
	if len(data) > 0 && data[len(data)-1] == "" {
		data = data[:len(data)-1]
	}
	if len(data) < 9 {
		ctx.JSON(400, gin.H{
			"error": "cannot parse response",
		})
		return
	}

	online, err := strconv.Atoi(data[4])
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot parse response",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}
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
