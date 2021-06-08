package server

import (
	"context"

	errors "innogrid.com/hcloud-classic/hcc_errors"
	rpcnovnc "innogrid.com/hcloud-classic/pb"

	"hcc/violin-novnc/action/grpc/errconv"
	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/model"
)

func (s *server) ControlVNC(ctx context.Context, in *rpcnovnc.ReqControlVNC) (*rpcnovnc.ResControlVNC, error) {
	var vncInfo model.Vnc
	var errStack *errors.HccErrorStack = nil
	var res rpcnovnc.ResControlVNC
	var result string

	vnc := in.GetVnc()
	vncInfo.ServerUUID = vnc.GetServerUUID()
	vncInfo.WebSocket = ""

	switch vnc.GetAction() {
	case "CREATE":
		errStack = driver.VNCD.Create(&vncInfo)
		if errStack != nil {
			res.HccErrorStack = errconv.HccStackToGrpc(errStack)
			result = "FAIL"
			goto EXIT
		}
		result = "Success"

	case "DELETE":
		errStack = driver.VNCD.Delete(&vncInfo)
		if errStack != nil {
			res.HccErrorStack = errconv.HccStackToGrpc(errStack)
			result = "FAIL"
			goto EXIT
		}
		result = "Success"

	case "UPDATE":
		fallthrough
	case "INFO":
		fallthrough
	default:
		vnc.Action = "UNDEFINED"
		logger.Logger.Println("Undefined Action: " + vnc.GetAction())
		errStack = errors.NewHccErrorStack(errors.NewHccError(
			errors.ViolinNoVNCGrpcOperationFail,
			"Undefined Action("+vnc.GetAction()+")"))
		res.HccErrorStack = errconv.HccStackToGrpc(errStack)
		result = "Fail"
		goto EXIT
	}

	res.Port = vncInfo.WebSocket

EXIT:
	_ = dao.InsertVNCRequestLog(vncInfo, vnc.GetUserID(), vnc.GetAction(), result)

	return &res, nil
}
