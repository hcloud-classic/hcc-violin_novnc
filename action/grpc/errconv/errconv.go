package errconv

import (
	errg "hcc/violin-novnc/action/grpc/pb/rpcmsgType"
	errh "hcc/violin-novnc/lib/errors"
)

func GrpcToHcc(eg *errg.HccError) *errh.HccError {
	return errh.NewHccError(eg.GetErrCode(), eg.GetErrText())
}

func HccToGrpc(eh *errh.HccError) *errg.HccError {
	return &errg.HccError{ErrCode: eh.ErrCode, ErrText: eh.ErrText}
}

func GrpcStackToHcc(esg *[]*errg.HccError) *errh.HccErrorStack {
	errStack := errh.NewHccErrorStack()

	for _, e := range *esg {
		errStack.Push(errh.NewHccError(e.GetErrCode(), e.GetErrText()))
	}

	return errStack
}

func HccStackToGrpc(esh *errh.HccErrorStack) *[]*errg.HccError {
	ges := []*errg.HccError{}
	for i := esh.Len(); i > 0; i-- {
		e := esh.Pop()
		ge := &errg.HccError{ErrCode: e.ErrCode, ErrText: e.ErrText}
		ges = append([]*errg.HccError{ge}, ges...)
	}
	return &ges
}
