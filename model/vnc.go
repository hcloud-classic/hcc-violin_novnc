package model

import errors "innogrid.com/hcloud-classic/hcc_errors"

// Vnc : Model strucy of vnc
type Vnc struct {
	ServerUUID string           `json:"server_uuid"`
	ServerIP   string           `json:"server_ip"`
	WebSocket  string           `json:"socket_number"`
	UserCount  string           `json:"user_cnt"`
	LastUserd  string           `json:"last_recently_used"`
	Errors     *errors.HccError `json:"error"`
}
