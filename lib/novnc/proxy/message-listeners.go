package proxy

import (
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/novnc/client"
	"hcc/violin-novnc/lib/novnc/common"
	"hcc/violin-novnc/lib/novnc/server"
)

type ClientUpdater struct {
	conn *client.ClientConn
}

// Consume recieves vnc-server-bound messages (Client messages) and updates the server part of the proxy
func (cc *ClientUpdater) Consume(seg *common.RfbSegment) error {
	//logger.Logger.Printf("ClientUpdater.Consume (vnc-server-bound): got segment type=%s bytes: %v", seg.SegmentType, seg.Bytes)
	switch seg.SegmentType {

	case common.SegmentFullyParsedClientMessage:
		clientMsg := seg.Message.(common.ClientMessage)
		//logger.Logger.Printf("ClientUpdater.Consume:(vnc-server-bound) got ClientMessage type=%s", clientMsg.Type())
		switch clientMsg.Type() {

		case common.SetPixelFormatMsgType:
			// update pixel format
			//logger.Logger.Println("ClientUpdater.Consume: updating pixel format")
			pixFmtMsg := clientMsg.(*server.MsgSetPixelFormat)
			cc.conn.PixelFormat = pixFmtMsg.PF
		}

		err := clientMsg.Write(cc.conn)
		if err != nil {
			logger.Logger.Printf("ClientUpdater.Consume (vnc-server-bound, SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err
	}
	return nil
}

type ServerUpdater struct {
	conn *server.ServerConn
}

func (p *ServerUpdater) Consume(seg *common.RfbSegment) error {

	//logger.Logger.Printf("WriteTo.Consume (ServerUpdater): got segment type=%s, object type:%d", seg.SegmentType, seg.UpcomingObjectType)
	switch seg.SegmentType {
	case common.SegmentMessageStart:
	case common.SegmentRectSeparator:
	case common.SegmentServerInitMessage:
		serverInitMessage := seg.Message.(*common.ServerInit)
		p.conn.SetHeight(serverInitMessage.FBHeight)
		p.conn.SetWidth(serverInitMessage.FBWidth)
		p.conn.SetDesktopName(string(serverInitMessage.NameText))
		p.conn.SetPixelFormat(&serverInitMessage.PixelFormat)

	case common.SegmentBytes:
		//logger.Logger.Printf("WriteTo.Consume (ServerUpdater SegmentBytes): got bytes len=%d", len(seg.Bytes))
		_, err := p.conn.Write(seg.Bytes)
		if err != nil {
			logger.Logger.Printf("WriteTo.Consume (ServerUpdater SegmentBytes): problem writing to port: %s", err)
		}
		return err
	case common.SegmentFullyParsedClientMessage:

		clientMsg := seg.Message.(common.ClientMessage)
		//logger.Logger.Printf("WriteTo.Consume (ServerUpdater): got ClientMessage type=%s", clientMsg.Type())
		err := clientMsg.Write(p.conn)
		if err != nil {
			logger.Logger.Printf("WriteTo.Consume (ServerUpdater SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err
	default:
		//return errors.New("WriteTo.Consume: undefined RfbSegment type")
	}
	return nil
}
