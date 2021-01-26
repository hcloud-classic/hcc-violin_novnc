package config

import (
	"github.com/Terry-Mao/goconf"
	errors "innogrid.com/hcloud-classic/hcc_errors"
)

var conf = goconf.New()
var config = vncConfig{}
var err error

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql config").Fatal()
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql id").Fatal()
	}

	Mysql.Password, err = config.MysqlConfig.String("password")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql password").Fatal()
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql address").Fatal()
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql port").Fatal()
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql database").Fatal()
	}
}

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "grpc config").Fatal()
	}

	Grpc = grpc{}
	Grpc.Port, err = config.GrpcConfig.Int("port")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "grpc port").Fatal()
	}

}

func parseHarp() {
	config.HarpConfig = conf.Get("harp")
	if config.HarpConfig == nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp config").Fatal()
	}

	Harp = harp{}
	Harp.Address, err = config.HarpConfig.String("harp_server_address")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp server address").Fatal()
	}
	Harp.Port, err = config.HarpConfig.Int("harp_server_port")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp server port").Fatal()
	}
	Harp.RequestTimeoutMs, err = config.HarpConfig.Int("harp_request_timeout_ms")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp timeout").Fatal()
	}
}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, err.Error()).Fatal()
	}

	parseMysql()
	parseGrpc()
	parseHarp()
}
