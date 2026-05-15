package mcp

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"simpleWebUtils/router/minecraft"
	"simpleWebUtils/utils"
)

type request struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type requestPayload struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      json.RawMessage        `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *errorField `json:"error,omitempty"`
}

type errorField struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"name":    "simpleWebUtils MCP Server",
		"version": "1.0.0",
		"endpoint": gin.H{
			"jsonrpc": "POST /mcp",
		},
		"methods": []string{
			"initialize",
			"tools/list",
			"tools/call",
		},
	})
}

func handleJSONRPC(ctx *gin.Context) {
	req, hasID, err := parseRequest(ctx)
	if err != nil {
		ctx.JSON(400, response{
			JSONRPC: "2.0",
			Error: &errorField{
				Code:    -32700,
				Message: "parse error",
			},
		})
		return
	}

	if req.JSONRPC != "2.0" {
		message := "invalid request: jsonrpc must be \"2.0\""
		if req.JSONRPC != "" {
			message = "invalid request: jsonrpc must be \"2.0\", got \"" + req.JSONRPC + "\""
		}
		ctx.JSON(200, response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &errorField{
				Code:    -32600,
				Message: message,
			},
		})
		return
	}

	result, rpcErr := invoke(req.Method, req.Params)
	if !hasID {
		ctx.Status(204)
		return
	}
	if rpcErr != nil {
		ctx.JSON(200, response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &errorField{
				Code:    rpcErr.Code,
				Message: rpcErr.Message,
			},
		})
		return
	}

	ctx.JSON(200, response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	})
}

type rpcError struct {
	Code    int
	Message string
}

func parseRequest(ctx *gin.Context) (request, bool, error) {
	var req request

	rawBody, err := ctx.GetRawData()
	if err != nil {
		return req, false, err
	}

	var payload requestPayload
	err = json.Unmarshal(rawBody, &payload)
	if err != nil {
		return req, false, err
	}

	req = request{
		JSONRPC: payload.JSONRPC,
		Method:  payload.Method,
		Params:  payload.Params,
	}
	if len(payload.ID) > 0 {
		err = json.Unmarshal(payload.ID, &req.ID)
		if err != nil {
			return req, false, err
		}
	}

	hasID := len(payload.ID) > 0
	if req.Params == nil {
		req.Params = map[string]interface{}{}
	}

	return req, hasID, nil
}

func invoke(method string, params map[string]interface{}) (interface{}, *rpcError) {
	switch method {
	case "initialize":
		return gin.H{
			"protocolVersion": "2024-11-05",
			"capabilities": gin.H{
				"tools": gin.H{},
			},
			"serverInfo": gin.H{
				"name":    "simpleWebUtils",
				"version": "1.0.0",
			},
		}, nil
	case "tools/list":
		return gin.H{
			"tools": []gin.H{
				{
					"name":        "bedrock_motd",
					"description": "Query Minecraft Bedrock MOTD",
					"inputSchema": gin.H{
						"type": "object",
						"properties": gin.H{
							"server": gin.H{
								"type":        "string",
								"description": "Minecraft Bedrock server address",
							},
							"port": gin.H{
								"type":        "string",
								"description": "Minecraft Bedrock port, default 19132",
							},
						},
						"required": []string{"server"},
					},
				},
				{
					"name":        "java_motd",
					"description": "Query Minecraft Java MOTD",
					"inputSchema": gin.H{
						"type": "object",
						"properties": gin.H{
							"server": gin.H{
								"type":        "string",
								"description": "Minecraft Java server address",
							},
							"port": gin.H{
								"type":        "string",
								"description": "Minecraft Java port, default 25565",
							},
						},
						"required": []string{"server"},
					},
				},
			},
		}, nil
	case "tools/call":
		name, ok := params["name"].(string)
		if !ok || name == "" {
			return nil, &rpcError{
				Code:    -32602,
				Message: "invalid params: name is required",
			}
		}

		args, ok := params["arguments"].(map[string]interface{})
		if !ok {
			args = map[string]interface{}{}
		}

		data, err := callTool(name, args)
		if err != nil {
			return nil, &rpcError{
				Code:    -32000,
				Message: utils.LocalAddressCleaner(err.Error()),
			}
		}

		serialized, err := json.Marshal(data)
		if err != nil {
			return nil, &rpcError{
				Code:    -32001,
				Message: "cannot serialize tool response",
			}
		}
		return gin.H{
			"content": []gin.H{
				{
					"type": "text",
					"text": string(serialized),
				},
			},
			"structuredContent": data,
		}, nil
	default:
		return nil, &rpcError{
			Code:    -32601,
			Message: "method not found",
		}
	}
}

func callTool(name string, arguments map[string]interface{}) (interface{}, error) {
	switch name {
	case "bedrock_motd":
		server, err := getRequiredString(arguments, "server")
		if err != nil {
			return nil, err
		}
		port, err := getOptionalString(arguments, "port")
		if err != nil {
			return nil, err
		}
		return minecraft.QueryBedrockMOTD(server, port)
	case "java_motd":
		server, err := getRequiredString(arguments, "server")
		if err != nil {
			return nil, err
		}
		port, err := getOptionalString(arguments, "port")
		if err != nil {
			return nil, err
		}
		return minecraft.QueryJavaMOTD(server, port)
	default:
		return nil, errors.New("unknown tool: " + name)
	}
}

func getRequiredString(values map[string]interface{}, key string) (string, error) {
	value, ok := values[key]
	if !ok {
		return "", errors.New("invalid params: " + key + " is required")
	}
	stringValue, ok := value.(string)
	if !ok || stringValue == "" {
		return "", errors.New("invalid params: " + key + " must be a non-empty string")
	}
	return stringValue, nil
}

func getOptionalString(values map[string]interface{}, key string) (string, error) {
	value, ok := values[key]
	if !ok {
		return "", nil
	}
	stringValue, ok := value.(string)
	if !ok {
		return "", errors.New("invalid params: " + key + " must be a string")
	}
	return stringValue, nil
}
