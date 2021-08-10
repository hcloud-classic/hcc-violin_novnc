package server

import (
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/novnc/common"
	"io"
	"net"
)

var DefaultClientMessages = []common.ClientMessage{
	&MsgSetPixelFormat{},
	&MsgSetEncodings{},
	&MsgFramebufferUpdateRequest{},
	&MsgKeyEvent{},
	&MsgPointerEvent{},
	&MsgClientCutText{},
	&MsgClientQemuExtendedKey{},
}

// FramebufferUpdate holds a FramebufferUpdate wire format message.
type FramebufferUpdate struct {
	_       [1]byte             // padding
	NumRect uint16              // number-of-rectangles
	Rects   []*common.Rectangle // rectangles
}

type ServerHandler func(*ServerConfig, *ServerConn) error

type ServerConfig struct {
	SecurityHandlers []SecurityHandler
	Encodings        []common.IEncoding
	PixelFormat      *common.PixelFormat
	ColorMap         *common.ColorMap
	ClientMessages   []common.ClientMessage
	DesktopName      []byte
	Height           uint16
	Width            uint16
	UseDummySession  bool

	//handler to allow for registering for messages, this can't be a channel
	//because of the websockets handler function which will kill the connection on exit if conn.handle() is run on another thread
	NewConnHandler ServerHandler
}

func wsHandlerFunc(ws io.ReadWriter, cfg *ServerConfig, sessionId string) {
	err := attachNewServerConn(ws, cfg, sessionId)
	if err != nil {
		logger.Logger.Println("Error attaching new connection. " + err.Error())
	}
}

func WsServe(url string, cfg *ServerConfig, server **WsServer) error {
	*server = &WsServer{cfg, nil, nil}
	(*server).Listen(url, WsHandler(wsHandlerFunc))
	return nil
}

func TcpServe(url string, cfg *ServerConfig) error {
	ln, err := net.Listen("tcp", url)
	if err != nil {
		logger.Logger.Println("Error listen. " + err.Error())
		return err
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			return err
		}
		go attachNewServerConn(c, cfg, "dummySession")
	}
	return nil
}

func attachNewServerConn(c io.ReadWriter, cfg *ServerConfig, sessionID string) error {

	conn, err := NewServerConn(c, cfg)
	if err != nil {
		return err
	}

	if err := ServerVersionHandler(cfg, conn); err != nil {
		logger.Logger.Printf("err: %v\n", err)
		_ = conn.Close()
		return err
	}

	if err := ServerSecurityHandler(cfg, conn); err != nil {
		_ = conn.Close()
		return err
	}

	//run the handler for this new incoming connection from a vnc-client
	//this is done before the init sequence to allow listening to server-init messages
	// (and maybe even interception in the future)
	err = cfg.NewConnHandler(cfg, conn)
	if err != nil {
		_ = conn.Close()
		return err
	}

	if err := ServerClientInitHandler(cfg, conn); err != nil {
		_ = conn.Close()
		return err
	}

	if err := ServerServerInitHandler(cfg, conn); err != nil {
		_ = conn.Close()
		return err
	}

	conn.SessionId = sessionID
	if cfg.UseDummySession {
		conn.SessionId = "dummySession"
	}

	//go here will kill ws connections
	_ = conn.handle()

	return nil
}
