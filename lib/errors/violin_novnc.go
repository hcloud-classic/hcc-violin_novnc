package errors

const (
	ViolinNoVNCInternalInitFail        uint64 = violinNoVNC + internal + initFail
	ViolinNoVNCInternalConnectionFail  uint64 = violinNoVNC + internal + connectionFail
	ViolinNoVNCInternalParsingError    uint64 = violinNoVNC + internal + parsingError
	ViolinNoVNCInternalDBOperationFail uint64 = violinNoVNC + internal + dbOperationFail

	ViolinNoVNCDriverRequestError       uint64 = violinNoVNC + driver + requestError
	ViolinNoVNCDriverResponseError      uint64 = violinNoVNC + driver + responseError
	ViolinNoVNCDriverReceiveError       uint64 = violinNoVNC + driver + receiveError
	ViolinNoVNCDriverParsingError       uint64 = violinNoVNC + driver + parsingError
	ViolinNoVNCDriverJsonUnmarshalError uint64 = violinNoVNC + driver + jsonUnmarshalError

	ViolinNoVNCGrpcArgumentError      uint64 = violinNoVNC + grpc + argumentError
	ViolinNoVNCGrpcRequestError       uint64 = violinNoVNC + grpc + requestError
	ViolinNoVNCGrpcSendError          uint64 = violinNoVNC + grpc + sendError
	ViolinNoVNCGrpcReceiveError       uint64 = violinNoVNC + grpc + receiveError
	ViolinNoVNCGrpcJsonUnmarshalError uint64 = violinNoVNC + grpc + jsonUnmarshalError
	ViolinNoVNCGrpcParsingError       uint64 = violinNoVNC + grpc + parsingError
	ViolinNoVNCGrpcConnectionFail     uint64 = violinNoVNC + grpc + connectionFail
	ViolinNoVNCGrpcUndefinedError     uint64 = violinNoVNC + grpc + undefinedError
	ViolinNoVNCGrpcDBOperationFail    uint64 = violinNoVNC + grpc + dbOperationFail
)
