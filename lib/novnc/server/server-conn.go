package server

import (
	"encoding/binary"
	"fmt"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/novnc/common"
	"io"
	"sync"
)

type ServerConn struct {
	c   io.ReadWriter
	cfg *ServerConfig

	protocol string
	m        sync.Mutex
	// If the pixel format uses a color map, then this is the color
	// map that is used. This should not be modified directly, since
	// the data comes from the server.
	colorMap *common.ColorMap

	// Name associated with the desktop, sent from the server.
	desktopName string

	// Encodings supported by the client. This should not be modified
	// directly. Instead, MsgSetEncodings() should be used.
	encodings []common.IEncoding

	// Height of the frame buffer in pixels, sent to the client.
	fbHeight uint16

	// Width of the frame buffer in pixels, sent to the client.
	fbWidth uint16

	// The pixel format associated with the connection. This shouldn't
	// be modified. If you wish to set a new pixel format, use the
	// SetPixelFormat method.
	pixelFormat *common.PixelFormat

	// a consumer for the parsed messages, to allow for recording and proxy
	Listeners *common.MultiListener

	SessionId string

	quit chan struct{}
}

// func (c *IServerConn) UnreadByte() error {
// 	return c.br.UnreadByte()
// }

func NewServerConn(c io.ReadWriter, cfg *ServerConfig) (*ServerConn, error) {
	// if cfg.ClientMessageCh == nil {
	// 	return nil, fmt.Errorf("ClientMessageCh nil")
	// }

	if len(cfg.ClientMessages) == 0 {
		return nil, fmt.Errorf("ClientMessage 0")
	}

	return &ServerConn{
		c: c,
		//br:          bufio.NewReader(c),
		//bw:          bufio.NewWriter(c),
		cfg:         cfg,
		quit:        make(chan struct{}),
		encodings:   cfg.Encodings,
		pixelFormat: cfg.PixelFormat,
		fbWidth:     cfg.Width,
		fbHeight:    cfg.Height,
		Listeners:   &common.MultiListener{},
	}, nil
}

func (c *ServerConn) Conn() io.ReadWriter {
	return c.c
}

func (c *ServerConn) SetEncodings(encs []common.EncodingType) error {
	encodings := make(map[int32]common.IEncoding)
	for _, enc := range c.cfg.Encodings {
		encodings[enc.Type()] = enc
	}
	for _, encType := range encs {
		if enc, ok := encodings[int32(encType)]; ok {
			c.encodings = append(c.encodings, enc)
		}
	}
	return nil
}

func (c *ServerConn) SetProtoVersion(pv string) {
	c.protocol = pv
}

func (c *ServerConn) Close() error {
	return c.c.(io.ReadWriteCloser).Close()
}

func (c *ServerConn) Read(buf []byte) (int, error) {
	return c.c.Read(buf)
}

func (c *ServerConn) Write(buf []byte) (int, error) {
	//	c.m.Lock()
	//	defer c.m.Unlock()
	return c.c.Write(buf)
}

func (c *ServerConn) ColorMap() *common.ColorMap {
	return c.colorMap
}

func (c *ServerConn) SetColorMap(cm *common.ColorMap) {
	c.colorMap = cm
}
func (c *ServerConn) DesktopName() string {
	return c.desktopName
}
func (c *ServerConn) CurrentPixelFormat() *common.PixelFormat {
	return c.pixelFormat
}
func (c *ServerConn) SetDesktopName(name string) {
	c.desktopName = name
}
func (c *ServerConn) SetPixelFormat(pf *common.PixelFormat) error {
	c.pixelFormat = pf
	return nil
}
func (c *ServerConn) Encodings() []common.IEncoding {
	return c.encodings
}
func (c *ServerConn) Width() uint16 {
	return c.fbWidth
}
func (c *ServerConn) Height() uint16 {
	return c.fbHeight
}
func (c *ServerConn) Protocol() string {
	return c.protocol
}
func (c *ServerConn) SetWidth(w uint16) {
	c.fbWidth = w
}
func (c *ServerConn) SetHeight(h uint16) {
	c.fbHeight = h
}

func (c *ServerConn) handle() error {

	defer func() {
		c.Listeners.Consume(&common.RfbSegment{
			SegmentType: common.SegmentConnectionClosed,
		})
	}()

	//create a map of all message types
	clientMessages := make(map[common.ClientMessageType]common.ClientMessage)
	for _, m := range c.cfg.ClientMessages {
		clientMessages[m.Type()] = m
	}

	for {
		select {
		case <-c.quit:
			return nil
		default:
			var messageType common.ClientMessageType
			if err := binary.Read(c, binary.BigEndian, &messageType); err != nil {
				logger.Logger.Printf("ServerConn.handle error: %v", err)
				return err
			}
			// logger.Logger.Printf("ServerConn.handle: got messagetype, %d", messageType)
			msg, ok := clientMessages[messageType]
			// logger.Logger.Printf("ServerConn.handle: found message type, %v", ok)
			if !ok {
				logger.Logger.Printf("ServerConn.handle: unsupported message-type: %v", messageType)
			}
			parsedMsg, err := msg.Read(c)
			// logger.Logger.Printf("ServerConn.handle: got parsed messagetype, %v", parsedMsg)
			//update connection for pixel format / color map changes
			switch parsedMsg.Type() {
			case common.SetPixelFormatMsgType:
				// update pixel format
				//logger.Logger.Println("ClientUpdater.Consume: updating pixel format")
				pixFmtMsg := parsedMsg.(*MsgSetPixelFormat)
				c.SetPixelFormat(&pixFmtMsg.PF)
				if pixFmtMsg.PF.TrueColor != 0 {
					c.SetColorMap(&common.ColorMap{})
				}
			}
			////////

			if err != nil {
				logger.Logger.Printf("srv err %s", err.Error())
				return err
			}

			//logger.Logger.Printf("IServerConn.Handle got ClientMessage: %s, %v", parsedMsg.Type(), parsedMsg)
			//TODO: treat set encodings by allowing only supported encoding in proxy configurations
			//// if parsedMsg.Type() == common.SetEncodingsMsgType{
			//// 	c.cfg.Encodings
			//// }

			seg := &common.RfbSegment{
				SegmentType: common.SegmentFullyParsedClientMessage,
				Message:     parsedMsg,
			}
			err = c.Listeners.Consume(seg)
			if err != nil {
				logger.Logger.Printf("IServerConn.Handle: listener consume err %s", err.Error())
				return err
			}
		}
	}
}
