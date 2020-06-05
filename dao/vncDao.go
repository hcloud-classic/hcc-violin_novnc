package dao

import (
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/mysql"
	"hcc/violin-novnc/model"
	"strconv"
	"time"
)

func CheckoutSpecificWSPort(WSPort string) (interface{}, error) {
	var IsPortAvailable string

	sql := "select if (isnull(server_uuid),'None','Using') from server_vnc where ws_port=?"

	err := mysql.Db.QueryRow(sql, WSPort).Scan(
		&IsPortAvailable)

	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	return IsPortAvailable, nil
}

func FindAvailableWsPort() (interface{}, error) {
	var serverUUID string
	var TargetIP string
	var TargetPort string
	var WebSocket string
	var TargetPass string
	var createdAt time.Time
	var AvailablePort string

	sql := "select * from server_vnc where ws_port=(select max(ws_port) from server_vnc) "
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&serverUUID, &TargetIP, &TargetPort, &WebSocket, &TargetPass, &createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		Port, parseerr := strconv.Atoi(WebSocket)
		if parseerr != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		AvailablePort = strconv.Itoa(Port + 1)
	}
	if AvailablePort == "" {
		AvailablePort = "5901"
	}
	return AvailablePort, nil

}

// CreateVNC : VNC DB create
func CreateVNC(args map[string]interface{}) (model.Vnc, error) {
	var serverVnc model.Vnc
	serverVnc = model.Vnc{
		ServerUUID:     args["server_uuid"].(string),
		TargetIP:       args["target_ip"].(string),
		TargetPort:     args["target_port"].(string),
		WebSocket:      args["websocket_port"].(string),
		TargetPass:     "qwe1212",
		ActionClassify: "Create",
	}
	sql := "insert into server_vnc(server_uuid, target_ip, target_port, ws_port, target_pass ,created_at) values (?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)

	if err != nil {
		logger.Logger.Println(err)
		return serverVnc, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	result, err := stmt.Exec(serverVnc.ServerUUID, serverVnc.TargetIP, serverVnc.TargetPort, serverVnc.WebSocket, serverVnc.TargetPass)

	if err != nil {
		logger.Logger.Println("DB Insert Error", err)
		return serverVnc, err
	}
	logger.Logger.Println(result.LastInsertId())
	serverVnc.Info = "Created"
	logger.Logger.Println("[Violin-novnc] Server VNC Create : ", result)
	return serverVnc, nil
}
