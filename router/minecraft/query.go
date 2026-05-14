package minecraft

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	// Bedrock unconnected ping payload used for MOTD queries.
	bedrockUnconnectedPingPayload = "0100000000240D12D300FFFF00FEFEFEFEFDFDFDFD12345678"
	maxVarIntBytes                = 5
)

type BedrockMOTDResult struct {
	Raw       []string `json:"raw"`
	Server    string   `json:"server"`
	ServerIP  string   `json:"server_ip"`
	Port      string   `json:"port"`
	MOTD      string   `json:"motd"`
	Protocol  string   `json:"protocol"`
	Version   string   `json:"version"`
	Online    int      `json:"online"`
	Max       int      `json:"max"`
	Level     string   `json:"level"`
	LevelName string   `json:"levelName"`
	Gamemode  string   `json:"gamemode"`
}

type JavaMOTDResult struct {
	Raw      map[string]interface{} `json:"raw"`
	Server   string                 `json:"server"`
	ServerIP string                 `json:"server_ip"`
	Port     string                 `json:"port"`
	MOTD     string                 `json:"motd"`
	Protocol int                    `json:"protocol"`
	Version  string                 `json:"version"`
	Online   int                    `json:"online"`
	Max      int                    `json:"max"`
}

func QueryBedrockMOTD(server string, port string) (*BedrockMOTDResult, error) {
	if server == "" {
		return nil, errors.New("server is required")
	}
	if port == "" {
		port = "19132"
	}

	addr, err := net.ResolveUDPAddr("udp", server+":"+port)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve server: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to server: %w", err)
	}
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	payload, err := hex.DecodeString(bedrockUnconnectedPingPayload)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare payload: %w", err)
	}
	_, err = conn.Write(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot send payload: %w", err)
	}

	buf := make([]byte, 2048)
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("cannot set read deadline: %w", err)
	}

	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot receive response: %w", err)
	}

	data, err := parseBedrockPayload(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	online, err := strconv.Atoi(data[4])
	if err != nil {
		return nil, fmt.Errorf("cannot parse online count: %w", err)
	}
	max, err := strconv.Atoi(data[5])
	if err != nil {
		return nil, fmt.Errorf("cannot parse max players: %w", err)
	}

	result := &BedrockMOTDResult{
		Raw:      data,
		Server:   server,
		ServerIP: addr.IP.String(),
		Port:     port,
		MOTD:     data[1],
		Protocol: data[2],
		Version:  data[3],
		Online:   online,
		Max:      max,
	}
	if len(data) > 6 {
		result.Level = data[6]
	}
	if len(data) > 7 {
		result.LevelName = data[7]
	}
	if len(data) > 8 {
		result.Gamemode = data[8]
	}

	return result, nil
}

func QueryJavaMOTD(server string, port string) (*JavaMOTDResult, error) {
	if server == "" {
		return nil, errors.New("server is required")
	}
	if port == "" {
		port = "25565"
	}

	addr, err := net.ResolveTCPAddr("tcp", server+":"+port)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve server: %w", err)
	}

	conn, err := net.DialTimeout("tcp", addr.String(), 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to server: %w", err)
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("cannot set deadline: %w", err)
	}

	handshakeData := make([]byte, 0)
	handshakeData = append(handshakeData, 0x00)
	versionVarInt, err := writeVarInt(47)
	if err != nil {
		return nil, fmt.Errorf("cannot build version varint: %w", err)
	}
	handshakeData = append(handshakeData, versionVarInt...)
	serverString, err := writeString(server)
	if err != nil {
		return nil, fmt.Errorf("cannot build server string: %w", err)
	}
	handshakeData = append(handshakeData, serverString...)
	portInt, err := strconv.Atoi(port)
	if err != nil || portInt < 0 || portInt > 65535 {
		return nil, fmt.Errorf("invalid port: %s", port)
	}
	handshakeData = append(handshakeData, byte(portInt>>8), byte(portInt))
	nextStateVarInt, err := writeVarInt(1)
	if err != nil {
		return nil, fmt.Errorf("cannot build next-state varint: %w", err)
	}
	handshakeData = append(handshakeData, nextStateVarInt...)

	packetLengthVarInt, err := writeVarInt(len(handshakeData))
	if err != nil {
		return nil, fmt.Errorf("cannot build packet-length varint: %w", err)
	}
	handshakePacket := append(packetLengthVarInt, handshakeData...)
	_, err = conn.Write(handshakePacket)
	if err != nil {
		return nil, fmt.Errorf("cannot send handshake payload: %w", err)
	}

	requestPacket := []byte{0x01, 0x00}
	_, err = conn.Write(requestPacket)
	if err != nil {
		return nil, fmt.Errorf("cannot send status request payload: %w", err)
	}

	packetData, err := readPacket(conn)
	if err != nil {
		return nil, fmt.Errorf("cannot read response packet: %w", err)
	}

	reader := strings.NewReader(string(packetData))
	packetID, err := readVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response packet id: %w", err)
	}
	if packetID != 0 {
		return nil, fmt.Errorf("unexpected response packet id: %d", packetID)
	}

	statusJSON, err := readString(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response payload: %w", err)
	}

	var raw map[string]interface{}
	err = json.Unmarshal([]byte(statusJSON), &raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse status json: %w", err)
	}

	versionName, protocol := extractVersion(raw["version"])
	online, max := extractPlayers(raw["players"])

	return &JavaMOTDResult{
		Raw:      raw,
		Server:   server,
		ServerIP: addr.IP.String(),
		Port:     port,
		MOTD:     flattenDescription(raw["description"]),
		Protocol: protocol,
		Version:  versionName,
		Online:   online,
		Max:      max,
	}, nil
}

func parseBedrockPayload(payload []byte) ([]string, error) {
	motd := string(payload)
	start := strings.Index(motd, "MCPE;")
	if start == -1 {
		return nil, errors.New("missing MCPE response marker")
	}
	motd = strings.TrimRight(motd[start:], "\x00\r\n\t ")
	data := strings.Split(motd, ";")
	if len(data) > 0 && data[len(data)-1] == "" {
		data = data[:len(data)-1]
	}
	if len(data) < 6 {
		return nil, errors.New("response fields are incomplete")
	}
	return data, nil
}

func readPacket(r io.Reader) ([]byte, error) {
	bufferedReader := bufio.NewReader(r)
	length, err := readVarInt(bufferedReader)
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return nil, errors.New("invalid packet length")
	}
	data := make([]byte, length)
	_, err = io.ReadFull(bufferedReader, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeVarInt(value int) ([]byte, error) {
	if value < 0 {
		return nil, errors.New("varint cannot be negative")
	}
	result := make([]byte, 0)
	for {
		if (value & ^0x7F) == 0 {
			result = append(result, byte(value))
			return result, nil
		}
		result = append(result, byte((value&0x7F)|0x80))
		value >>= 7
	}
}

func readVarInt(r io.ByteReader) (int, error) {
	numRead := 0
	result := 0
	for {
		read, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		value := int(read & 0x7F)
		result |= value << (7 * numRead)
		numRead++
		if numRead > maxVarIntBytes {
			return 0, errors.New("varint is too big")
		}
		if (read & 0x80) == 0 {
			break
		}
	}
	return result, nil
}

func writeString(value string) ([]byte, error) {
	strBytes := []byte(value)
	lengthVarInt, err := writeVarInt(len(strBytes))
	if err != nil {
		return nil, err
	}
	return append(lengthVarInt, strBytes...), nil
}

func readString(r io.ByteReader) (string, error) {
	length, err := readVarInt(r)
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", errors.New("string length cannot be negative")
	}
	strBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		strBytes[i], err = r.ReadByte()
		if err != nil {
			return "", err
		}
	}
	return string(strBytes), nil
}

func extractVersion(value interface{}) (string, int) {
	versionMap, ok := value.(map[string]interface{})
	if !ok {
		return "", 0
	}

	name, _ := versionMap["name"].(string)
	protocolFloat, _ := versionMap["protocol"].(float64)

	return name, int(protocolFloat)
}

func extractPlayers(value interface{}) (int, int) {
	playersMap, ok := value.(map[string]interface{})
	if !ok {
		return 0, 0
	}

	onlineFloat, _ := playersMap["online"].(float64)
	maxFloat, _ := playersMap["max"].(float64)

	return int(onlineFloat), int(maxFloat)
}

func flattenDescription(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return typed
	case map[string]interface{}:
		text, _ := typed["text"].(string)
		result := text
		if extraValue, ok := typed["extra"].([]interface{}); ok {
			for _, item := range extraValue {
				result += flattenDescription(item)
			}
		}
		if result != "" {
			return result
		}
		translate, _ := typed["translate"].(string)
		return translate
	default:
		return ""
	}
}
