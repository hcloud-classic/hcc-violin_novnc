package model

import errors "innogrid.com/hcloud-classic/hcc_errors"

// Vnc : Model struct of vnc
type Vnc struct {
	ServerUUID string           `json:"server_uuid"`
	ServerIP   string           `json:"server_ip"`
	WebSocket  string           `json:"port_number"`
	UserCount string           `json:"user_cnt"`
	LastUsed  string           `json:"last_recently_used"`
	Errors    *errors.HccError `json:"error"`
}
