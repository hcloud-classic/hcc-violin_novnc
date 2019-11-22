package config

import (
	"fmt"
	"hcc/violin-novnc/lib/logger"

	"github.com/Terry-Mao/goconf"
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

	HTTP.Port, err = config.HTTPConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

}

func parseVnc() {
	config.VncConfig = conf.Get("violin_novnc")
	if config.VncConfig == nil {
		logger.Logger.Panicln("no violin_novnc section")
	}

	ViolinNovnc = violin_novnc{}
	ViolinNovnc.ServerAddress, err = config.VncConfig.String("violin_novnc_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	fmt.Println(ViolinNovnc.ServerAddress)
	ViolinNovnc.ServerPort, err = config.VncConfig.Int("violin_novnc_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	ViolinNovnc.RequestTimeoutMs, err = config.VncConfig.Int("violin_novnc_request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}
func parseRabbitMQ() {
	config.RabbitMQConfig = conf.Get("rabbitmq")
	if config.RabbitMQConfig == nil {
		logger.Logger.Panicln("no rabbitmq section")
	}

	RabbitMQ = rabbitmq{}
	RabbitMQ.ID, err = config.RabbitMQConfig.String("rabbitmq_id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Password, err = config.RabbitMQConfig.String("rabbitmq_password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Address, err = config.RabbitMQConfig.String("rabbitmq_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Port, err = config.RabbitMQConfig.Int("rabbitmq_port")
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
	parseVnc()
}
