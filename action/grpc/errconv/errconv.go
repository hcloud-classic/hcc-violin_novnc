package errconv

import (
	errh "github.com/hcloud-classic/hcc_errors"
	errg "hcc/violin-novnc/action/grpc/pb/rpcmsgType"
)

func GrpcToHcc(eg *errg.HccError) *errh.HccError {
	return errh.NewHccError(eg.GetErrCode(), eg.GetErrText())
}

func HccToGrpc(eh *errh.HccError) *errg.HccError {
	return &errg.HccError{ErrCode: eh.Code(), ErrText: eh.Text()}
}

func GrpcStackToHcc(esg *[]*errg.HccError) *errh.HccErrorStack {
	errStack := errh.NewHccErrorStack()

	for _, e := range *esg {
		errStack.Push(errh.NewHccError(e.GetErrCode(), e.GetErrText()))
	}

	hccErrStack := *errStack
	es := hccErrStack[1:]
	return &es
}

func HccStackToGrpc(esh *errh.HccErrorStack) []*errg.HccError {
	ges := []*errg.HccError{}
	for i := 0; i <= esh.Len(); i++ {
		ge := &errg.HccError{ErrCode: (*esh)[i].Code(), ErrText: (*esh)[i].Text()}
		ges = append(ges, ge)
	}
	return ges
}
