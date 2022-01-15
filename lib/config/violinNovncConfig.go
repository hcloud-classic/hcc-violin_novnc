package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/violin-novnc/violin-novnc.conf"

type vncConfig struct {
	RsakeyConfig *goconf.Section
	MysqlConfig  *goconf.Section
	GrpcConfig   *goconf.Section
	HornConfig   *goconf.Section
	HarpConfig   *goconf.Section
}
