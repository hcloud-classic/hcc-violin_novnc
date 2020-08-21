package player

import (
	"encoding/binary"

	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/novnc/client"
	"hcc/violin-novnc/lib/novnc/common"
	"hcc/violin-novnc/lib/novnc/server"
	"io"
	"time"
)

type VncStreamFileReader interface {
	io.Reader
	CurrentTimestamp() int
	ReadStartSession() (*common.ServerInit, error)
	CurrentPixelFormat() *common.PixelFormat
	Encodings() []common.IEncoding
}

type FBSPlayListener struct {
	Conn             *server.ServerConn
	Fbs              VncStreamFileReader
	serverMessageMap map[uint8]common.ServerMessage
	firstSegDone     bool
	startTime        int
}

func ConnectFbsFile(filename string, conn *server.ServerConn) (*FbsReader, error) {
	fbs, err := NewFbsReader(filename)
	if err != nil {
		logger.Logger.Println("failed to open fbs reader:", err)
		return nil, err
	}
	//NewFbsReader("/Users/amitbet/vncRec/recording.rbs")
	initMsg, err := fbs.ReadStartSession()
	if err != nil {
		logger.Logger.Println("failed to open read fbs start session:", err)
		return nil, err
	}
	conn.SetPixelFormat(&initMsg.PixelFormat)
	conn.SetHeight(initMsg.FBHeight)
	conn.SetWidth(initMsg.FBWidth)
	conn.SetDesktopName(string(initMsg.NameText))

	return fbs, nil
}

func NewFBSPlayListener(conn *server.ServerConn, r *FbsReader) *FBSPlayListener {
	h := &FBSPlayListener{Conn: conn, Fbs: r}
	cm := client.MsgBell(0)
	h.serverMessageMap = make(map[uint8]common.ServerMessage)
	h.serverMessageMap[0] = &client.MsgFramebufferUpdate{}
	h.serverMessageMap[1] = &client.MsgSetColorMapEntries{}
	h.serverMessageMap[2] = &cm
	h.serverMessageMap[3] = &client.MsgServerCutText{}

	return h
}
func (handler *FBSPlayListener) Consume(seg *common.RfbSegment) error {

	switch seg.SegmentType {
	case common.SegmentFullyParsedClientMessage:
		clientMsg := seg.Message.(common.ClientMessage)
		//logger.Logger.Printf("ClientUpdater.Consume:(vnc-server-bound) got ClientMessage type=%s", clientMsg.Type())
		switch clientMsg.Type() {

		case common.FramebufferUpdateRequestMsgType:
			if !handler.firstSegDone {
				handler.firstSegDone = true
				handler.startTime = int(time.Now().UnixNano() / int64(time.Millisecond))
			}
			handler.sendFbsMessage()
		}
		// server.MsgFramebufferUpdateRequest:
	}
	return nil
}

func (h *FBSPlayListener) sendFbsMessage() {
	var messageType uint8
	//messages := make(map[uint8]common.ServerMessage)
	fbs := h.Fbs
	//conn := h.Conn
	err := binary.Read(fbs, binary.BigEndian, &messageType)
	if err != nil {
		logger.Logger.Println("TestServer.NewConnHandler: Error in reading FBS segment: ", err)
		return
	}
	//common.IClientConn{}
	binary.Write(h.Conn, binary.BigEndian, messageType)
	msg := h.serverMessageMap[messageType]
	if msg == nil {
		logger.Logger.Println("TestServer.NewConnHandler: Error unknown message type: ", messageType)
		return
	}
	timeSinceStart := int(time.Now().UnixNano()/int64(time.Millisecond)) - h.startTime
	timeToSleep := fbs.CurrentTimestamp() - timeSinceStart
	if timeToSleep > 0 {
		time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	}

	err = msg.CopyTo(fbs, h.Conn, fbs)
	if err != nil {
		logger.Logger.Println("TestServer.NewConnHandler: Error in reading FBS segment: ", err)
		return
	}
}
