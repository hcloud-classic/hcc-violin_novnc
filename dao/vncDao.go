package dao

import (
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/mysql"
	"hcc/violin-novnc/model"
)

// CreateVNC : VNC DB createS
func CreateVNC(args map[string]interface{}) (model.Vnc, error) {
	// fmt.Println("@@@params@@@\n", args["server_uuid"].(string), args["target_ip"].(string), args["target_port"].(string), args["websocket_port"].(string))
	// fmt.Println("allocWsPort: ", allocWsPort)
	var serverVnc model.Vnc
	serverVnc = model.Vnc{
		ServerUUID:     args["server_uuid"].(string),
		TargetIP:       args["target_ip"].(string),
		TargetPort:     args["target_port"].(string),
		WebSocket:      args["websocket_port"].(string),
		TargetPass:     "qwe1212",
		ActionClassify: "Create",
	}
	sql := "INSERT INTO server_vnc(server_uuid, target_ip, target_port, ws_port, target_pass ,created_at) values (?, ?, ?, ?, ?, now()) ON DUPLICATE KEY UPDATE target_ip=?, target_port=?, ws_port=?, created_at=now()"
	stmt, err := mysql.Db.Prepare(sql)

	if err != nil {
		logger.Logger.Println(err)
		return serverVnc, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	result, err := stmt.Exec(serverVnc.ServerUUID, serverVnc.TargetIP, serverVnc.TargetPort, serverVnc.WebSocket, serverVnc.TargetPass, serverVnc.TargetIP, serverVnc.TargetPort, serverVnc.WebSocket)

	if err != nil {
		logger.Logger.Println("DB Insert Error", err)
		return serverVnc, err
	}

	serverVnc.Info = "Created"
	logger.Logger.Println("[Violin-novnc] Server VNC Create : ", result)
	return serverVnc, nil
}

func DeleteVNC(srvUUID string) error {
	sql := "DELETE FROM `violin_novnc`.`server_vnc` WHERE server_uuid=\"" + srvUUID + "\""

	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	return nil
}

func GetVNCServerList() ([]string, error) {
	var srvUUIDList []string
	sql := "SELECT `server_uuid` FROM `violin_novnc`.`server_vnc`"

	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		var uuid string
		err := stmt.Scan(&uuid)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		srvUUIDList = append(srvUUIDList, uuid)
	}
	return srvUUIDList, nil
}

func InitVNCServer() error {
	sql := "TRUNCATE TABLE `violin_novnc`.`server_vnc`"

	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	return nil

}
