package driver

import (
	"sync"

	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/lib/errors"
	"hcc/violin-novnc/lib/logger"
	vncproxy "hcc/violin-novnc/lib/novnc/proxy"
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

func (vncd *VNCDriver) Prepare() {
	err := dao.InitVNCServer()
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}
}

func (vncd *VNCDriver) Create(srvUUID string) (string, *errors.HccErrorStack) {
	var srvIP, port string
	var es *errors.HccErrorStack = nil

	logger.Logger.Print("Find exist VNC proxy websocket...")
	wsPort, ok := vncd.serverWSMap.Load(srvUUID)
	if !ok {
		logger.Logger.Println("[FAIL]")
		logger.Logger.Print("Asking server ip to harp...")

		srvIP, es = client.RC.GetServerIP(srvUUID)
		if es != nil {
			logger.Logger.Println("[FAIL]")
			es.Push(errors.NewHccError(errors.ViolinNoVNCDriverReceiveError, "GetServerIP"))
			return "", es
		}
		logger.Logger.Println("[SUCCESS] -- ", srvIP)

		logger.Logger.Print("Find available port...")
		port = PD.GetAvailablePort()
		if port == "0" {
			logger.Logger.Println("[FAIL]")
			es.Push(errors.NewHccError(errors.ViolinNoVNCDriverReceiveError, "GetAvailablePort"))
			return "", es
		}
		logger.Logger.Println("[SUCCESS] -- ", port)

		vncd.serverWSMap.Store(srvUUID, port)

		wsURL := "http://0.0.0.0:" + port + "/" + srvUUID + "_" + port

		proxy := &vncproxy.VncProxy{
			WsListeningURL: wsURL,
			SingleSession: &vncproxy.VncSession{
				Target:         srvIP,
				TargetPort:     vncd.vncPort,
				TargetPassword: vncd.vncPasswd, //"vncPass",
				ID:             srvUUID,
				Status:         vncproxy.SessionStatusInit,
				Type:           vncproxy.SessionTypeProxyPass,
			}, // to be used when not using sessions
			UsingSessions: false, //false = single session - defined in the var above
		}
		vncd.serverProxyMap.Store(srvUUID, proxy)

		logger.Logger.Println("[SUCCESS]")
		var args = make(map[string]interface{})
		args["server_uuid"] = srvUUID
		args["target_ip"] = srvIP
		args["target_port"] = vncd.vncPort
		args["websocket_port"] = port

		_, err := dao.CreateVNC(args)
		if err != nil {
			logger.Logger.Println(err.Error())
			es.Push(err)
			return "", es
		}

		vncd.serverConnectionMap.Store(srvUUID, 1)

		go func() {
			logger.Logger.Print("Create VNC Proxy...")

			err := proxy.StartListening()
			if err != nil {
				logger.Logger.Println("[FAIL]\n", err)
			}
			vncd.serverConnectionMap.Delete(srvUUID)

			p, _ := vncd.serverWSMap.Load(srvUUID)
			PD.ReturnPort(p.(string))
			vncd.serverWSMap.Delete(srvUUID)
			logger.Logger.Println(srvUUID, " proxy Server Successfully Closed")
		}()

		return port, nil
	}
	logger.Logger.Println("[SUCCESS] -- " + port)

	if cn, b := vncd.serverConnectionMap.Load(srvUUID); b {
		vncd.serverConnectionMap.Store(srvUUID, cn.(int)+1)
	}

	return wsPort.(string), nil
}

func (vncd *VNCDriver) Delete(srvUUID string) *errors.HccErrorStack {
	var es *errors.HccErrorStack = nil

	if cn, b := vncd.serverConnectionMap.Load(srvUUID); b {
		if cn.(int) > 1 {
			vncd.serverConnectionMap.Store(srvUUID, cn.(int)-1)
		} else {
			// stop vnc proxy server
			if proxy, b := vncd.serverProxyMap.Load(srvUUID); b {
				logger.Logger.Println(srvUUID, " Proxy will close")
				proxy.(*vncproxy.VncProxy).Shutdown()
				if err := dao.DeleteVNC(srvUUID); err != nil {
					es = errors.NewHccErrorStack(err)
				}
			}
		}
	}
	return es
}
