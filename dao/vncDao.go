package dao

import (
	"hcc/violin-novnc/lib/mysql"
	"hcc/violin-novnc/model"

	errors "innogrid.com/hcloud-classic/hcc_errors"
)

func sendStmt(sql string, params ...interface{}) (mysql.Result, *errors.HccError) {
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		return nil, errors.NewHccError(errors.ViolinNoVNCInternalOperationFail, "sql Prepare : "+err.Error())
	}

	defer func() {
		_ = stmt.Close()
	}()

	result, err := stmt.Exec(params...)

	if err != nil {
		return result, errors.NewHccError(errors.ViolinNoVNCInternalOperationFail, "stmt Exec : "+err.Error())
	}

	return result, nil
}

func sendQuery(sql string) (*mysql.Rows, *errors.HccError) {
	result, err := mysql.Db.Query(sql)
	if err != nil {
		return nil, errors.NewHccError(errors.ViolinNoVNCInternalOperationFail, "sql Query : "+err.Error())
	}

	return result, nil
}

// CreateVNC : VNC DB createS
func InsertVNCInfo(vncInfo model.Vnc) *errors.HccError {

	sql := "INSERT IGNORE INTO `violin_novnc`.`allocated_socket`(socket_number, server_uuid, user_cnt, last_recently_used) values (?, ?, ?, now());"

	_, err := sendStmt(sql, vncInfo.WebSocket, vncInfo.ServerUUID, vncInfo.UserCount)
	if err != nil {
		return err
	}

	return nil
}

func DeleteVNCInfo(vncInfo model.Vnc) *errors.HccError {
	sql := "DELETE FROM `violin_novnc`.`allocated_socket` WHERE socket_number = ?"

	_, err := sendStmt(sql, vncInfo.WebSocket)
	if err != nil {
		return err
	}

	return nil
}

func IncreaseVNCUserCount(vncInfo model.Vnc) *errors.HccError {

	sql := "UPDATE `violin_novnc`.`allocated_socket` SET user_cnt = user_cnt + 1 WHERE socket_number = ?;"

	_, err := sendStmt(sql, vncInfo.WebSocket)
	if err != nil {
		return err
	}

	return nil
}

func DecreaseVNCUserCount(vncInfo model.Vnc) *errors.HccError {

	sql := "UPDATE `violin_novnc`.`allocated_socket` SET user_cnt = user_cnt - 1 WHERE socket_number = ?;"

	_, err := sendStmt(sql, vncInfo.WebSocket)
	if err != nil {
		return err
	}

	return nil
}

func GetVNCSrvSockPair() (*mysql.Rows, *errors.HccError) {

	sql := "SELECT `socket_number`, `server_uuid`, `user_cnt` FROM `violin_novnc`.`allocated_socket` ORDER BY `socket_number` ASC"

	result, err := sendQuery(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func InsertVNCRequestLog(vncInfo model.Vnc, userID, action, result string) *errors.HccError {
	sql := "INSERT INTO `violin_novnc`.`vnc_connection_log`(server_uuid, user, request_type, result) values (?, ?, ?, ?);"

	_, err := sendStmt(sql, vncInfo.ServerUUID, userID, action, result)
	if err != nil {
		return err
	}

	return nil
}
