package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/violin-novnc/violin-novnc.conf"

type vncConfig struct {
	MysqlConfig    *goconf.Section
	HTTPConfig     *goconf.Section
	RabbitMQConfig *goconf.Section
	VncConfig      *goconf.Section
}

/*-----------------------------------
         Config File Example

##### CONFIG START #####

[http]
port 7800

[violin_novnc]
violin_novnc_server_address 10.0.100.9
violin_novnc_server_port 7800
violin_novnc_request_timeout_ms 5000


-----------------------------------*/
