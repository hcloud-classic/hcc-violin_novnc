package driver

import (
	"errors"
	"flag"
	"fmt"
	"github.com/graphql-go/graphql"
	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/model"
	vncproxy "hcc/violin-novnc/proxy"
	"os"
	"sync"
)

//**node scheduling argument */
// cpu, mem, end of bmc ip address

//RunProcxy :RunProcxy
func RunProcxy(params graphql.ResolveParams, wsURL string) {
	//create default session if required
	// recorddir string, target string, targPass string, wsport string
	var recordDir string
	var targetVnc string
	var targetVncPass string
	var wsPort string
	recordDir = "/var/log/violin-novnc/recordings/" + params.Args["server_uuid"].(string)
	targetVnc = params.Args["target_ip"].(string) + ":" + params.Args["target_port"].(string)
	targetVncPass = params.Args["target_pass"].(string)
	wsPort = params.Args["websocket_port"].(string)
	//Not use
	fmt.Println(recordDir, "   ", targetVnc, "    ", targetVncPass, "    ", wsPort)

	err := logger.CreateDirIfNotExist("/var/log/violin-novnc/recordings/")
	if err != nil {
		logger.Logger.Println(err)
		return
	}

	err = logger.CreateDirIfNotExist(recordDir)
	if err != nil {
		logger.Logger.Println(err)
		return
	}

	var vncPass string
	var targetVncPort string
	var targetVncHost string
	var tcpPort string
	// var wsPort = flag.String("wsPort", "", "websocket port")
	// var targetVncPass = flag.String("targPass", "", "target vnc password")
	// var recordDir = flag.String("recDir", "", "path to save FBS recordings WILL NOT RECORD if not defined.")
	// var targetVnc = flag.String("target", "", "target vnc server (host:port or /path/to/unix.socket)")
	// var tcpPort = flag.String("tcpPort", "", "tcp port")

	// var vncPass = flag.String("vncPass", "", "password on incoming vnc connections to the proxy, defaults to no password")
	// var targetVncPort = flag.String("targPort", "", "target vnc server port (deprecated, use -target)")
	// var targetVncHost = flag.String("targHost", "", "target vnc server host (deprecated, use -target)")
	// var logLevel = flag.String("logLevel", "info", "change logging level")

	if tcpPort == "" && wsPort == "" {
		logger.Logger.Println("no listening port defined")
		flag.Usage()
		os.Exit(1)
	}

	if targetVnc == "" && targetVncPort == "" {
		flag.Usage()
		logger.Logger.Println("no target vnc server host/port or socket defined")
		return
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
}

var mutex = &sync.Mutex{}

func Runner(params graphql.ResolveParams) (interface{}, error) {
	vnc := model.Vnc{
		ActionClassify: params.Args["action"].(string),
	}

	var err error
	if params.Args["action"].(string) != "" {
		allocWsPort, errs := dao.FindAvailableWsPort()
		if errs != nil {
			vnc.Info = "Web Socket Not found"
			return vnc, nil
		} else {
			vnc.WebSocket = allocWsPort.(string)
		}

		// fmt.Println("allocWsPort: ", allocWsPort, "result.WebSocket ", params.Args["websocket_port"])
		switch params.Args["action"].(string) {
		case "Create":
			params.Args["websocket_port"] = vnc.WebSocket
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

		wsURL := "http://0.0.0.0:" + vnc.WebSocket + "/" + params.Args["server_uuid"].(string) + "_" + vnc.WebSocket
		go func() {
			mutex.Lock()
			RunProcxy(params, wsURL)
			mutex.Unlock()
		}()
	}

	return vnc, nil
}
