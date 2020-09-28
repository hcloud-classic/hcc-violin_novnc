package config

type harp struct {
	Address          string `goconf:"harp:harp_server_address"`
	Port             int64  `goconf:"harp:harp_server_port"`
	RequestTimeoutMs int64  `goconf:"http:harp_request_timeout_ms"`
}

var Harp harp
