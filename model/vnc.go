package model

// Vnc : Model strucy of vnc
type Vnc struct {
	ServerUUID     string `json:"server_uuid"`
	TargetIP       string `json:"target_ip"`
	TargetPort     string `json:"target_port"`
	WebSocket      string `json:"websocket_port"`
	WsURL          string `json:"ws_url"`
	TargetPass     string `json:"target_pass"`
	Info           string `json:"vnc_info"`
	ActionClassify string `json:"action"`
}
