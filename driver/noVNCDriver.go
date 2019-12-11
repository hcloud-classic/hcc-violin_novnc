package driver

import (
	"errors"
	"flag"
	"fmt"
	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/model"
	vncproxy "hcc/violin-novnc/proxy"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/graphql-go/graphql"
)

//**node scheduling argument */
// cpu, mem, end of bmc ip address

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
	//recordDir := "/var/log/violin-novnc/recordings/" + params.Args["server_uuid"].(string) + "_" + wsPort
	recordDir := ""
	//Not use
	fmt.Println(recordDir, "   ", targetVnc, "    ", targetVncPass, "    ", wsPort)

	//err := logger.CreateDirIfNotExist("/var/log/violin-novnc/recordings/")
	//if err != nil {
	//	logger.Logger.Println(err)
	//	return err
	//}

	// err = logger.CreateDirIfNotExist(recordDir)
	// if err != nil {
	// 	logger.Logger.Println(err)
	// 	return err
	// }

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
		for {
			genWsPort = GenerateRandPort(40000, 50000)
			if SelfCheckPortScan(genWsPort) == "Closed" {
				break
			}
		}

		allocWsPort, errs := dao.CheckoutSpecificWSPort(strconv.Itoa(genWsPort))
		// allocWsPort, errs := dao.CheckoutSpecificWSPort("59245")
		// fmt.Println("allocWsPort : ", allocWsPort)
		// fmt.Println("errs : ", errs)
		if errs == nil {
			vnc.Info = "Web Socket Not found"
			return vnc, nil
		} else {
			if allocWsPort == nil {
				vnc.WebSocket = strconv.Itoa(genWsPort)

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
