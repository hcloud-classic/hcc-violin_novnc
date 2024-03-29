package config

type grpc struct {
	Port                 int64 `goconf:"grpc:port"`                    // Port : Port number for listening graphql request via grpc server
	ClientPingIntervalMs int64 `goconf:"grpc:client_ping_interval_ms"` // ClientPingIntervalMs : Interval for pinging gRPC servers
	ClientPingTimeoutMs  int64 `goconf:"grpc:client_ping_timeout_ms"`  // ClientPingTimeoutMs : Timeout for pinging gRPC servers
}

// Grpc : grpc config structure
var Grpc grpc
