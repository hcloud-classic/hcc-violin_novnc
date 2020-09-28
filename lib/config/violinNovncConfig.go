package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/violin-novnc/violin-novnc.conf"

type vncConfig struct {
	MysqlConfig *goconf.Section
	HTTPConfig  *goconf.Section
	HarpConfig  *goconf.Section
}
