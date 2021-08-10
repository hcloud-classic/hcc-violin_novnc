package config

import (
	"github.com/Terry-Mao/goconf"
	"innogrid.com/hcloud-classic/hcc_errors"
)

var conf = goconf.New()
var config = vncConfig{}
var err error

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, "no mysql section").Fatal()
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Mysql.Password, err = config.MysqlConfig.String("password")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}
	Mysql.ConnectionRetryCount, err = config.MysqlConfig.Int("connection_retry_count")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Mysql.ConnectionRetryIntervalMs, err = config.MysqlConfig.Int("connection_retry_interval_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

}

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, "no grpc section").Fatal()
	}

	Grpc.Port, err = config.GrpcConfig.Int("port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Grpc.ClientPingIntervalMs, err = config.GrpcConfig.Int("client_ping_interval_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Grpc.ClientPingTimeoutMs, err = config.GrpcConfig.Int("client_ping_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}
}

func parseHarp() {
	config.HarpConfig = conf.Get("harp")
	if config.HarpConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, "no harp section").Fatal()
	}

	Harp = harp{}
	Harp.ServerAddress, err = config.HarpConfig.String("harp_server_address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Harp.ServerPort, err = config.HarpConfig.Int("harp_server_port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}

	Harp.RequestTimeoutMs, err = config.HarpConfig.Int("harp_request_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, err.Error()).Fatal()
	}
}

// Init : Parse config file and initialize config structure
func Init() {
	if err = conf.Parse(configLocation); err != nil {
		hcc_errors.NewHccError(hcc_errors.ViolinNoVNCInternalParsingError, err.Error()).Fatal()
	}

	parseMysql()
	parseGrpc()
	parseHarp()
}
