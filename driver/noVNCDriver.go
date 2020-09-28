package driver

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/graphql-go/graphql"

	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/driver/grpccli"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/model"
	vncproxy "hcc/violin-novnc/proxy"
)

type VNCManager struct {
	serverIPMap map[string]string /* [ServerUUID] ServerIP */
	serverWSMap map[string]string /* [ServerIP] Websocket */
	vncPort     string
	vncPasswd   string
}

var VNCM = VNCManager{
	serverIPMap: make(map[string]string),
	serverWSMap: make(map[string]string),
	vncPort:     "5901",
	vncPasswd:   "qwe1212",
}

func (vncm *VNCManager) Create(token, srvUUID string) (string, error) {
	var srvIP string
	var err error
	logger.Logger.Println("Find server ip...")
	srvIP, ok := vncm.serverIPMap[srvUUID]
	if !ok {
		logger.Logger.Println("Asking server ip to harp")
		srvIP, err = grpccli.RC.GetServerIP(srvUUID)
		if err != nil {
			logger.Logger.Println(err)
			return "", err
		}
		vncm.serverIPMap[srvUUID] = srvIP
	}
	port := PM.GetAvailablePort()
	if port == "0" {
		return "", errors.New("No more websocket port available")
	}

	wsURL := "http://0.0.0.0:" + port + "/" + srvUUID + "_" + port

	proxy := &vncproxy.VncProxy{
		WsListeningURL: wsURL, // empty = not listening on ws
		SingleSession: &vncproxy.VncSession{
			Target:         srvIP + ":" + port,
			TargetPort:     vncm.vncPort,
			TargetPassword: vncm.vncPasswd, //"vncPass",
			ID:             "dummySession" + port,
			Status:         vncproxy.SessionStatusInit,
			Type:           vncproxy.SessionTypeProxyPass,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	go proxy.StartListening()

	return port, nil
}

//RunProcxy :RunProcxy
func RunProcxy(params graphql.ResolveParams) error {
	//create default session if required
	// recorddir string, target string, targPass string, wsport string
	var targetVnc string
	var targetVncPass string
	var wsPort string
	targetVnc = params.Args["target_ip"].(string) + ":" + params.Args["target_port"].(string)
	targetVncPass = params.Args["target_pass"].(string)
	wsPort = params.Args["websocket_port"].(string)
	wsURL := "http://0.0.0.0:" + wsPort + "/" + params.Args["server_uuid"].(string) + "_" + wsPort
	recordDir := ""
	//Not use
	fmt.Println(recordDir, "   ", targetVnc, "    ", targetVncPass, "    ", wsPort)

	var vncPass string
	var targetVncPort string
	var targetVncHost string
	var tcpPort string

	if tcpPort == "" && wsPort == "" {
		logger.Logger.Println("no listening port defined")
		flag.Usage()
		os.Exit(1)
	}

	if targetVnc == "" && targetVncPort == "" {
		flag.Usage()
		err := errors.New("no target vnc server host/port or socket defined")
		logger.Logger.Println(err)
		return err
	}

	if vncPass == "" {
		logger.Logger.Println("proxy will have no password")
	}

	tcpURL := ""
	if tcpPort != "" {
		tcpURL = ":" + string(tcpPort)
	}

	proxy := &vncproxy.VncProxy{
		WsListeningURL:   wsURL, // empty = not listening on ws
		TCPListeningURL:  tcpURL,
		ProxyVncPassword: vncPass, //empty = no auth
		SingleSession: &vncproxy.VncSession{
			Target:         targetVnc,
			TargetHostname: targetVncHost,
			TargetPort:     targetVncPort,
			TargetPassword: targetVncPass, //"vncPass",
			ID:             "dummySession",
			Status:         vncproxy.SessionStatusInit,
			Type:           vncproxy.SessionTypeProxyPass,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	if recordDir != "" {
		fullPath := recordDir
		// if err != nil {
		// 	logger.Error("bad recording path: ", err)
		// }
		logger.Logger.Println("FBS recording is turned on, writing to dir: ", fullPath)
		proxy.RecordingDir = fullPath
		proxy.SingleSession.Type = vncproxy.SessionTypeRecordingProxy
	} else {
		logger.Logger.Println("FBS recording is turned off")
	}
	proxy.StartListening()

	return nil
}

var mutex = &sync.Mutex{}

func Runner(params graphql.ResolveParams) (interface{}, error) {
	vnc := model.Vnc{
		ActionClassify: params.Args["action"].(string),
	}
	var err error
	if params.Args["action"].(string) != "" {
		mutex.Lock()
		var genWsPort int
		// allocWsPort, errs := dao.FindAvailableWsPort()
		//
	retry:
		for {
			genWsPort = GenerateRandPort(40000, 50000)
			if SelfCheckPortScan(genWsPort) == "Closed" {
				break
			}
		}

		//fmt.Println("genWsPort : ", genWsPort)
		allocWsPort, errs := dao.CheckoutSpecificWSPort(strconv.Itoa(genWsPort))
		// allocWsPort, errs := dao.CheckoutSpecificWSPort("59245")
		//fmt.Println("allocWsPort : ", allocWsPort)
		// fmt.Println("errs : ", errs)
		if errs != nil {
			vnc.Info = "Web Socket Not found"
			return vnc, nil
		} else {
			if allocWsPort.(string) == "None" {
				vnc.WebSocket = strconv.Itoa(genWsPort)
				fmt.Println("genWsPort:", genWsPort)
			} else {
				goto retry
			}
		}
		params.Args["websocket_port"] = vnc.WebSocket

		// fmt.Println("allocWsPort: ", allocWsPort, "result.WebSocket ", params.Args["websocket_port"])
		switch params.Args["action"].(string) {
		case "Create":
			vnc, err = dao.CreateVNC(params.Args)
			if err != nil {
				return err, nil
			}
		case "Delete":
		case "Update":
		case "Info":
		default:
			vnc.Info = "Failed Please Choose Action"
			return vnc, errors.New("[Violin-Novnc] : Please Choose Action")
		}

		// TODO
		// Need to close server when exit popup

		go func(params graphql.ResolveParams) {
			fmt.Println("websocket_port", params.Args["websocket_port"])
			var p = params
			_ = RunProcxy(p)
		}(params)

		mutex.Unlock()
	}

	return vnc, nil
}

func GenerateRandPort(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func SelfCheckPortScan(Port int) string {
	port := strconv.FormatInt(int64(Port), 10)
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	fmt.Println(conn)
	if err == nil {
		fmt.Println("Port", Port, "open")
		conn.Close()
		return "Open"
	}
	return "Closed"
}
