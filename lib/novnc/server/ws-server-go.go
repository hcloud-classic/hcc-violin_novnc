package server

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"

	"hcc/violin-novnc/lib/logger"
)

type WsServer struct {
	cfg         *ServerConfig
	proxyServer *http.Server
	mux         *http.ServeMux
}

type WsHandler func(io.ReadWriter, *ServerConfig, string)

func (wsServer *WsServer) Listen(urlStr string, handlerFunc WsHandler) {

	if urlStr == "" {
		urlStr = "/"
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		logger.Logger.Println("error while parsing url: ", err)
	}

	wsServer.mux = http.NewServeMux()

	wsServer.mux.Handle(url.Path, websocket.Handler(
		func(ws *websocket.Conn) {
			path := ws.Request().URL.Path
			var sessionId string
			if path != "" {
				sessionId = path[1:]
			}

			ws.PayloadType = websocket.BinaryFrame
			handlerFunc(ws, wsServer.cfg, sessionId)
		}))

	// err = http.ListenAndServe(url.Host, nil)
	wsServer.proxyServer = &http.Server{Addr: url.Host, Handler: wsServer.mux}
	logger.Logger.Println("Server start")
	err = wsServer.proxyServer.ListenAndServe()
	if err != nil {
		logger.Logger.Println("ListenAndServe: " + err.Error())
	}
}

func (wsServer *WsServer) Shutdown() {
	wsServer.proxyServer.Shutdown(context.Background())
}
