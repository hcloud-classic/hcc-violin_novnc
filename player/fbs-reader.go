package player

import (
	"bytes"
	"encoding/binary"
	"hcc/violin-novnc/common"
	"hcc/violin-novnc/encodings"
	"hcc/violin-novnc/lib/logger"
	"io"
	"os"
)

type FbsReader struct {
	reader           io.Reader
	buffer           bytes.Buffer
	currentTimestamp int
	pixelFormat      *common.PixelFormat
	encodings        []common.IEncoding
}

func (fbs *FbsReader) CurrentTimestamp() int {
	return fbs.currentTimestamp
}

func (fbs *FbsReader) Read(p []byte) (n int, err error) {
	if fbs.buffer.Len() < len(p) {
		seg, err := fbs.ReadSegment()

		if err != nil {
			logger.Logger.Println("FBSReader.Read: error reading FBSsegment: ", err)
			return 0, err
		}
		fbs.buffer.Write(seg.bytes)
		fbs.currentTimestamp = int(seg.timestamp)
	}
	return fbs.buffer.Read(p)
}

func (fbs *FbsReader) CurrentPixelFormat() *common.PixelFormat { return fbs.pixelFormat }

//func (fbs *FbsReader) CurrentColorMap() *common.ColorMap       { return &common.ColorMap{} }
func (fbs *FbsReader) Encodings() []common.IEncoding { return fbs.encodings }

func NewFbsReader(fbsFile string) (*FbsReader, error) {

	reader, err := os.OpenFile(fbsFile, os.O_RDONLY, 0644)
	if err != nil {
		logger.Logger.Println("NewFbsReader: can't open fbs file: ", fbsFile)
		return nil, err
	}
	return &FbsReader{reader: reader,
		encodings: []common.IEncoding{
			&encodings.CopyRectEncoding{},
			&encodings.ZLibEncoding{},
			&encodings.ZRLEEncoding{},
			&encodings.CoRREEncoding{},
			&encodings.HextileEncoding{},
			&encodings.TightEncoding{},
			&encodings.TightPngEncoding{},
			&encodings.EncCursorPseudo{},
			&encodings.EncLedStatePseudo{},
			&encodings.RawEncoding{},
			&encodings.RREEncoding{},
		},
	}, nil

}

func (fbs *FbsReader) ReadStartSession() (*common.ServerInit, error) {

	initMsg := common.ServerInit{}
	reader := fbs.reader

	var framebufferWidth uint16
	var framebufferHeight uint16
	var SecTypeNone uint32
	//read rfb header information (the only part done without the [size|data|timestamp] block wrapper)
	//.("FBS 001.000\n")
	bytes := make([]byte, 12)
	_, err := reader.Read(bytes)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init message - FBS file Version:", err)
		return nil, err
	}

	//read the version message into the buffer so it will be written in the first rbs block
	//RFB 003.008\n
	bytes = make([]byte, 12)
	_, err = fbs.Read(bytes)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init - RFB Version: ", err)
		return nil, err
	}

	//push sec type and fb dimensions
	binary.Read(fbs, binary.BigEndian, &SecTypeNone)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init - SecType: ", err)
	}

	//read frame buffer width, height
	binary.Read(fbs, binary.BigEndian, &framebufferWidth)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init - FBWidth: ", err)
		return nil, err
	}
	initMsg.FBWidth = framebufferWidth

	binary.Read(fbs, binary.BigEndian, &framebufferHeight)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init - FBHeight: ", err)
		return nil, err
	}
	initMsg.FBHeight = framebufferHeight

	//read pixel format
	pixelFormat := &common.PixelFormat{}
	binary.Read(fbs, binary.BigEndian, pixelFormat)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init - Pixelformat: ", err)
		return nil, err
	}
	initMsg.PixelFormat = *pixelFormat
	//read padding
	bytes = make([]byte, 3)
	fbs.Read(bytes)
	fbs.pixelFormat = pixelFormat

	//read desktop name
	var desknameLen uint32
	binary.Read(fbs, binary.BigEndian, &desknameLen)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init - deskname Len: ", err)
		return nil, err
	}
	initMsg.NameLength = desknameLen

	bytes = make([]byte, desknameLen)
	fbs.Read(bytes)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: error reading rbs init - desktopName: ", err)
		return nil, err
	}

	initMsg.NameText = bytes

	return &initMsg, nil
}

func (fbs *FbsReader) ReadSegment() (*FbsSegment, error) {
	reader := fbs.reader
	var bytesLen uint32

	//read length
	err := binary.Read(reader, binary.BigEndian, &bytesLen)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadStartSession: read len, error reading rbs file: ", err)
		return nil, err
	}

	paddedSize := (bytesLen + 3) & 0x7FFFFFFC

	//read bytes
	bytes := make([]byte, paddedSize)
	_, err = reader.Read(bytes)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadSegment: read bytes, error reading rbs file: ", err)
		return nil, err
	}

	//remove padding
	actualBytes := bytes[:bytesLen]

	//read timestamp
	var timeSinceStart uint32
	binary.Read(reader, binary.BigEndian, &timeSinceStart)
	if err != nil {
		logger.Logger.Println("FbsReader.ReadSegment: read timestamp, error reading rbs file: ", err)
		return nil, err
	}

	//timeStamp := time.Unix(timeSinceStart, 0)
	seg := &FbsSegment{bytes: actualBytes, timestamp: timeSinceStart}
	return seg, nil
}

type FbsSegment struct {
	bytes     []byte
	timestamp uint32
}
