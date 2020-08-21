package client

import (
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/novnc/common"
	"io"
)

type WriteTo struct {
	Writer io.Writer
	Name   string
}

func (p *WriteTo) Consume(seg *common.RfbSegment) error {

	//logger.Logger.Printf("WriteTo.Consume ("+p.Name+"): got segment type=%s", seg.SegmentType)
	switch seg.SegmentType {
	case common.SegmentMessageStart:
	case common.SegmentRectSeparator:
	case common.SegmentBytes:
		_, err := p.Writer.Write(seg.Bytes)
		if err != nil {
			logger.Logger.Printf("WriteTo.Consume ("+p.Name+" SegmentBytes): problem writing to port: %s", err)
		}
		return err
	case common.SegmentFullyParsedClientMessage:

		clientMsg := seg.Message.(common.ClientMessage)
		//logger.Logger.Printf("WriteTo.Consume ("+p.Name+"): got ClientMessage type=%s", clientMsg.Type())
		err := clientMsg.Write(p.Writer)
		if err != nil {
			logger.Logger.Printf("WriteTo.Consume ("+p.Name+" SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err
	default:
		//return errors.New("WriteTo.Consume: undefined RfbSegment type")
	}
	return nil
}
