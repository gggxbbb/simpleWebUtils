package code

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

var code_detail = map[int]string{

	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",

	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",

	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",

	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "Request-URI Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	420: "Enhance Your Calm",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	425: "Reserved for WebDAV advanced collections expired proposal",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	444: "No Response",
	450: "Blocked by Windows Parental Controls",
	451: "Unavailable For Legal Reasons",
	497: "HTTP Request Sent to HTTPS Port",
	498: "Invalid Token",
	499: "Client Closed Request",

	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	509: "Bandwidth Limit Exceeded",
	510: "Not Extended",
	511: "Network Authentication Required",
	521: "Web Server Is Down",
	522: "Connection Timed Out",
	523: "Origin Is Unreachable",
	525: "SSL Handshake Failed",
	530: "Site is frozen",
	598: "Network read timeout error",
}

func codeDetail(c *gin.Context) {
	code := c.Param("code")
	// parse code to int
	intCode, err := strconv.Atoi(code)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid code",
		})
		return
	}
	if code_detail[intCode] == "" {
		c.JSON(400, gin.H{
			"error": "invalid code",
			"see":   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status",
		})
		return
	}
	c.JSON(intCode, gin.H{
		"code":        intCode,
		"description": code_detail[intCode],
		"see":         "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/" + code,
		"cat":         "https://http.cat/" + code,
	})

}
