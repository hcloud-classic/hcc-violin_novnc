package server

import (
	"hcc/violin-novnc/lib/logger"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"
)

type WsServer struct {
	cfg *ServerConfig
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

	http.Handle(url.Path, websocket.Handler(
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
	err = http.ListenAndServe(url.Host, nil)
	if err != nil {
		logger.Logger.Println("ListenAndServe: " + err.Error())
	}
}
