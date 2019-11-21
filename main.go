package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/amitbet/vncproxy/logger"
	vncproxy "github.com/amitbet/vncproxy/proxy"
)

func init() {
	err := vncInit.MainInit()
	if err != nil {
		panic(err)
	}
}
func main() {
	RunProcxy("/var/log/violin-novnc/recordings/", "172.18.0.1:5901", "qwe1212", "5905")

	defer func() {
		schedulerEnd.MainEnd()
	}()

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
		return
	}
}

//RunProcxy :RunProcxy
func RunProcxy(recorddir string, target string, targPass string, wsport string) {
	//create default session if required
	var recordDir string
	var targetVnc string
	var targetVncPass string
	var wsPort string
	recordDir = recorddir
	targetVnc = target
	targetVncPass = targPass
	wsPort = wsport
	//Not use
	var vncPass string
	var targetVncPort string
	var targetVncHost string
	var logLevel string
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

	flag.Parse()
	logger.SetLogLevel(logLevel)

	if tcpPort == "" && wsPort == "" {
		logger.Error("no listening port defined")
		flag.Usage()
		os.Exit(1)
	}

	if targetVnc == "" && targetVncPort == "" {
		logger.Error("no target vnc server host/port or socket defined")
		flag.Usage()
		os.Exit(1)
	}

	if vncPass == "" {
		logger.Warn("proxy will have no password")
	}

	tcpURL := ""
	if tcpPort != "" {
		tcpURL = ":" + string(tcpPort)
	}
	wsURL := ""
	if wsPort != "" {
		wsURL = "http://0.0.0.0:" + string(wsPort) + "/"
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
		logger.Info("FBS recording is turned on, writing to dir: ", fullPath)
		proxy.RecordingDir = fullPath
		proxy.SingleSession.Type = vncproxy.SessionTypeRecordingProxy
	} else {
		logger.Info("FBS recording is turned off")
	}

	proxy.StartListening()
}
