package driver

import (
	"errors"
	"sync"

	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/driver/grpccli"
	"hcc/violin-novnc/lib/logger"
	vncproxy "hcc/violin-novnc/proxy"
)

type VNCManager struct {
	serverWSMap         sync.Map /* [ServerUUID(string)] ServerWS(string) */
	serverConnectionMap sync.Map /* [ServerUUID(string)] Connection Number(int) */
	createMutex         sync.Mutex
	addMutex            sync.Mutex
	vncPort             string
	vncPasswd           string
}

var VNCM = VNCManager{
	serverWSMap:         sync.Map{},
	serverConnectionMap: sync.Map{},
	createMutex:         sync.Mutex{},
	addMutex:            sync.Mutex{},
	vncPort:             "5901",
	vncPasswd:           "qwe1212",
}

func (vncm *VNCManager) Prepare() {
	/*
		srvUUIDList, err := dao.GetVNCServerList()
		if err != nil {
			logger.Logger.Fatalf("Fail to perparing server websocket")
		}
		for _, uuid := range srvUUIDList {
			vncm.Create("", uuid)
		}
	*/
	err := dao.InitVNCUser()
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}
}

func (vncm *VNCManager) Create(token, srvUUID string) (string, error) {
	var srvIP, port string
	var err error
	logger.Logger.Println("Find VNC proxy websocket")
	vncm.createMutex.Lock()
	wsPort, ok := vncm.serverWSMap.Load(srvUUID)
	if !ok {
		logger.Logger.Println("Asking server ip to harp")

		srvIP, err = grpccli.RC.GetServerIP(srvUUID)
		if err != nil {
			logger.Logger.Println(err)
			return "", err
		}

		port = PM.GetAvailablePort()
		if port == "0" {
			return "", errors.New("No more websocket port available")
		}

		vncm.serverWSMap.Store(srvUUID, port)
		vncm.createMutex.Unlock()
		vncm.addMutex.Lock() // Block user count add before proxy create

		wsURL := "http://0.0.0.0:" + port + "/" + srvUUID + "_" + port

		proxy := &vncproxy.VncProxy{
			WsListeningURL: wsURL,
			SingleSession: &vncproxy.VncSession{
				Target:         srvIP,
				TargetPort:     vncm.vncPort,
				TargetPassword: vncm.vncPasswd, //"vncPass",
				ID:             srvUUID,
				Status:         vncproxy.SessionStatusInit,
				Type:           vncproxy.SessionTypeProxyPass,
			}, // to be used when not using sessions
			UsingSessions: false, //false = single session - defined in the var above
		}

		go proxy.StartListening()

		var args = make(map[string]interface{})
		args["server_uuid"] = srvUUID
		args["target_ip"] = srvIP
		args["target_port"] = vncm.vncPort
		args["websocket_port"] = port

		_, err = dao.CreateVNC(args)
		if err != nil {
			logger.Logger.Println(err.Error())
			return "", err
		}

		vncm.serverConnectionMap.Store(srvUUID, 1)
		vncm.addMutex.Unlock()

		err = dao.AddVNCUser(token, srvUUID)
		if err != nil {
			logger.Logger.Println(err.Error())
			return "", err
		}
		return port, nil
	}
	vncm.createMutex.Unlock()
	logger.Logger.Println("WSPort Already exist " + port)

	vncm.addMutex.Lock()
	cn, _ := vncm.serverConnectionMap.Load(srvUUID)
	vncm.serverConnectionMap.Store(srvUUID, cn.(int)+1)
	vncm.addMutex.Unlock()

	err = dao.AddVNCUser(token, srvUUID)
	if err != nil {
		logger.Logger.Println(err.Error())
		return "", err
	}

	return wsPort.(string), nil
}

func (vncm *VNCManager) Delete(token, srvUUID string) error {
	return nil
}
