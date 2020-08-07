package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/violin-novnc/lib/logger"
)

var conf = goconf.New()
var config = vncConfig{}
var err error

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		logger.Logger.Panicln("no mysql section")
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Password, err = config.MysqlConfig.String("password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseHTTP() {
	config.HTTPConfig = conf.Get("http")
	if config.HTTPConfig == nil {
		logger.Logger.Panicln("no http section")
	}

	HTTP = http{}
	HTTP.Port, err = config.HTTPConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

}

func parseHarp() {
	config.HarpConfig = conf.Get("harp")
	if config.HarpConfig == nil {
		logger.Logger.Panicln("no harp section")
	}

	Harp = harp{}
	Harp.Address, err = config.HarpConfig.String("harp_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	Harp.Port, err = config.HarpConfig.Int("harp_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	Harp.RequestTimeoutMs, err = config.HarpConfig.Int("harp_request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseMysql()
	parseHTTP()
	parseHarp()
}
