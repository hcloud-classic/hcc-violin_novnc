package driver

import (
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"sync"

	errors "innogrid.com/hcloud-classic/hcc_errors"

	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/lib/logger"
	vncproxy "hcc/violin-novnc/lib/novnc/proxy"
	"hcc/violin-novnc/model"
)

type VNCDriver struct {
	serverWSMap         sync.Map /* [ServerUUID(string)] ServerWS(string) */
	serverConnectionMap sync.Map /* [ServerUUID(string)] Connection Number(int) */
	serverProxyMap      sync.Map /* [ServerUUID(string)] Proxy Server(*http.Server) */
	createMutex         sync.Mutex
	addMutex            sync.Mutex
	vncPort             string
	vncPasswd           string
}

var VNCD = VNCDriver{
	serverWSMap:         sync.Map{},
	serverConnectionMap: sync.Map{},
	serverProxyMap:      sync.Map{},
	createMutex:         sync.Mutex{},
	addMutex:            sync.Mutex{},
	vncPort:             "5901",
	vncPasswd:           "qwe1212",
}

func Init() *errors.HccErrorStack {
	var esInit = errors.NewHccErrorStack(
		errors.NewHccError(errors.ViolinNoVNCInternalInitFail, "noVNC Driver"))

	result, err := dao.GetVNCSrvSockPair()
	if err != nil {
		_ = esInit.Push(err)
		return esInit
	}

	defer func() {
		_ = result.Close()
	}()

	for result.Next() {
		var webSocket, srvUUID, userCount string

		if e := result.Scan(&webSocket, &srvUUID, &userCount); e != nil {
			_ = esInit.Push(errors.NewHccError(errors.ViolinNoVNCInternalOperationFail,
				"Fail to read socket-uuid pair data"))
			return esInit
		}

		vncInfo := model.Vnc{
			ServerUUID: srvUUID,
			WebSocket:  webSocket,
			UserCount:  userCount,
		}
		es := VNCD.Recover(&vncInfo)
		if es != nil {
			esInit.Merge(es)
			_ = dao.DeleteVNCInfo(vncInfo)
		}
	}

	if esInit.Len() > 1 { // esInit has one error by default
		return esInit
	}

	return nil
}

func (vncd *VNCDriver) setVncInfo(vncInfo *model.Vnc) *errors.HccErrorStack {

	var esSetInfo = errors.NewHccErrorStack(
		errors.NewHccError(errors.ViolinNoVNCDriverOperationFail, "Set VNC Info"))

	logger.Logger.Print("Asking server ip to harp...")

	srvIP, es := client.RC.GetServerIP(vncInfo.ServerUUID)
	if es != nil {
		logger.Logger.Println("[FAIL]")
		_ = esSetInfo.Push(errors.NewHccError(
			errors.ViolinNoVNCDriverReceiveError,
			"Cannot find server ip ["+vncInfo.ServerUUID+"]"))
		esSetInfo.Merge(es)

		return esSetInfo
	}
	vncInfo.ServerIP = srvIP
	logger.Logger.Println("[SUCCESS] -- ", vncInfo.ServerIP)

	if vncInfo.WebSocket == "" {
		logger.Logger.Print("Find available port...")

		// Port driver is not thread safe, use lock
		vncd.createMutex.Lock()
		vncInfo.WebSocket = PD.GetAvailablePort()

		if vncInfo.WebSocket == "" {
			logger.Logger.Println("[FAIL]")
			_ = esSetInfo.Push(errors.NewHccError(
				errors.ViolinNoVNCDriverOperationFail,
				"Websocket allocation for ["+vncInfo.ServerUUID+"]"))
			vncd.createMutex.Unlock()

			return esSetInfo
		}

		vncd.createMutex.Unlock()
		logger.Logger.Println("[SUCCESS] -- ", vncInfo.WebSocket)
	}
	vncd.serverWSMap.Store(vncInfo.ServerUUID, vncInfo.WebSocket)

	vncInfo.Errors = dao.InsertVNCInfo(*vncInfo)
	if vncInfo.Errors != nil {
		/**
		 * VNC info insert fail is not fatal error.
		 * We won't propagate this error to client.
		 * Just log it.
		**/
		logger.Logger.Print(vncInfo.Errors.Error())
	}

	return nil
}

func (vncd *VNCDriver) getSingleSessionProxy(vncInfo *model.Vnc) *vncproxy.VncProxy {

	proxy := vncproxy.VncProxy{
		WsListeningURL: "http://0.0.0.0:" + vncInfo.WebSocket + "/" + vncInfo.ServerUUID + "_" + vncInfo.WebSocket,
		SingleSession: &vncproxy.VncSession{
			Target:         vncInfo.ServerIP,
			TargetPort:     vncd.vncPort,
			TargetPassword: vncd.vncPasswd, //"vncPass",
			ID:             vncInfo.ServerUUID,
			// Status:         vncproxy.SessionStatusInit,
			Type: vncproxy.SessionTypeProxyPass,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	return &proxy
}

func (vncd *VNCDriver) Create(vncInfo *model.Vnc) *errors.HccErrorStack {
	var esCreate = errors.NewHccErrorStack(
		errors.NewHccError(errors.ViolinNoVNCDriverOperationFail, "Create VNC Proxy"))

	logger.Logger.Print("Find exist VNC proxy websocket...")
	wsPort, ok := vncd.serverWSMap.Load(vncInfo.ServerUUID)
	if !ok {
		logger.Logger.Println("[FAIL]")

		es := vncd.setVncInfo(vncInfo)
		if es != nil {
			esCreate.Merge(es)

			return esCreate
		}

		port, _ := strconv.Atoi(vncInfo.WebSocket)
		_, err := client.RC.CreatePortForwarding(&pb.ReqCreatePortForwarding{
			PortForwarding: &pb.PortForwarding{
				ServerUUID:   "master",
				ForwardTCP:   true,
				ForwardUDP:   false,
				ExternalPort: int64(port),
				InternalPort: 0,
				Description:  "VNC_" + vncInfo.ServerUUID,
			},
		})
		if err != nil {
			esCreate.Merge(es)

			return esCreate
		}

		proxy := vncd.getSingleSessionProxy(vncInfo)

		vncd.serverProxyMap.Store(vncInfo.ServerUUID, proxy)
		vncd.serverConnectionMap.Store(vncInfo.ServerUUID, 1)

		// run proxy
		go func() {
			defer func() { // clean up
				vncd.serverConnectionMap.Delete(vncInfo.ServerUUID)
				vncd.serverProxyMap.Delete(vncInfo.ServerUUID)

				p, _ := vncd.serverWSMap.Load(vncInfo.ServerUUID)
				PD.ReturnPort(p.(string))
				vncd.serverWSMap.Delete(vncInfo.ServerUUID)

				err := dao.DeleteVNCInfo(*vncInfo)
				if err != nil {
					// Not fatal error. Just log it
					logger.Logger.Println("Delete proxy server from database failed")
				}
			}()

			logger.Logger.Print("Create VNC Proxy...")

			err := proxy.StartListening()
			if err != nil {
				logger.Logger.Println("[FAIL] Start VNC proxy\n", err)
			}
		}()

		return nil
	}

	vncInfo.WebSocket = wsPort.(string)
	logger.Logger.Println("[SUCCESS] -- " + vncInfo.WebSocket)

	// Prevent load before store from another go rutine
	vncd.addMutex.Lock()
	if cn, b := vncd.serverConnectionMap.Load(vncInfo.ServerUUID); b {
		vncd.serverConnectionMap.Store(vncInfo.ServerUUID, cn.(int)+1)
	}
	vncd.addMutex.Unlock()

	err := dao.IncreaseVNCUserCount(*vncInfo)
	if err != nil {
		// Not fatal error. Just log it
		logger.Logger.Println("Proxy user count increase failed\n", err.Error())
	}

	return nil
}

func (vncd *VNCDriver) Delete(vncInfo *model.Vnc) *errors.HccErrorStack {
	var es *errors.HccErrorStack

	ws, b := vncd.serverWSMap.Load(vncInfo.ServerUUID)
	if !b {
		es = errors.NewHccErrorStack(
			errors.NewHccError(
				errors.ViolinNoVNCInternalOperationFail,
				"Cannot find Server"))
		return es
	}

	vncInfo.WebSocket = ws.(string)

	port, _ := strconv.Atoi(vncInfo.WebSocket)
	_, err := client.RC.DeletePortForwarding(&pb.ReqDeletePortForwarding{
		PortForwarding: &pb.PortForwarding{
			ServerUUID:   "master",
			ExternalPort: int64(port),
		},
	})
	if err != nil {
		es = errors.NewHccErrorStack(
			errors.NewHccError(
				errors.ViolinNoVNCGrpcRequestError,
				err.Error()))
		return es
	}

	// Prevent load before store in another go rutine
	vncd.addMutex.Lock()
	cn, _ := vncd.serverConnectionMap.Load(vncInfo.ServerUUID)
	if cn.(int) > 1 {
		vncd.serverConnectionMap.Store(vncInfo.ServerUUID, cn.(int)-1)
	} else {
		// stop vnc proxy server
		proxy, _ := vncd.serverProxyMap.Load(vncInfo.ServerUUID)
		logger.Logger.Println(vncInfo.ServerUUID, " Proxy will close")
		proxy.(*vncproxy.VncProxy).Shutdown()
	}

	vncd.addMutex.Unlock()

	err = dao.DecreaseVNCUserCount(*vncInfo)
	if err != nil {
		// Not fatal error. Just log it
		logger.Logger.Println("Proxy user count decrease failed\n", err.Error())
	}

	return nil
}

func (vncd *VNCDriver) Recover(vncInfo *model.Vnc) *errors.HccErrorStack {
	var esRestore = errors.NewHccErrorStack(
		errors.NewHccError(errors.ViolinNoVNCDriverOperationFail, "Restore VNC Proxy"))

	es := vncd.setVncInfo(vncInfo)
	if es != nil {
		esRestore.Merge(es)

		return esRestore
	}

	PD.SetLastPort(vncInfo.WebSocket)

	proxy := vncd.getSingleSessionProxy(vncInfo)

	vncd.serverProxyMap.Store(vncInfo.ServerUUID, proxy)
	userCnt, _ := strconv.Atoi(vncInfo.UserCount)
	vncd.serverConnectionMap.Store(vncInfo.ServerUUID, userCnt)

	// run proxy
	go func() {
		defer func() { // clean up
			vncd.serverConnectionMap.Delete(vncInfo.ServerUUID)
			vncd.serverProxyMap.Delete(vncInfo.ServerUUID)

			p, _ := vncd.serverWSMap.Load(vncInfo.ServerUUID)
			PD.ReturnPort(p.(string))
			vncd.serverWSMap.Delete(vncInfo.ServerUUID)

			_ = dao.DeleteVNCInfo(*vncInfo)
		}()

		logger.Logger.Print("Create VNC Proxy...")

		err := proxy.StartListening()
		if err != nil {
			logger.Logger.Println("[FAIL] Start VNC proxy\n", err)
		}
	}()
	return nil
}
