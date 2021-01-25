package config

type grpc struct {
	Port             int64 `goconf:"grpc:port"`               // Port : Port number for listening graphql request via grpc server
	RequestTimeoutMs int64 `goconf:"grpc:request_timeout_ms"` // RequestTimeoutMs : Timeout for Grpc request
}

// Grpc : grpc config structure
var Grpc grpc
