package server

import (
	"fmt"

	"github.com/dragonflyoss/Dragonfly/pkg/constants"
	"github.com/dragonflyoss/Dragonfly/pkg/errortypes"
)

// ResultInfo identify a struct that will returned to the client.
type ResultInfo struct {
	code int
	msg  string
	data interface{}
}

// NewResultInfoWithError returns a new ResultInfo with error only.
// And it will fill the result code according to the type of error.
func NewResultInfoWithError(err error) ResultInfo {
	if errortypes.IsEmptyValue(err) ||
		errortypes.IsInvalidValue(err) {
		return NewResultInfoWithCodeError(constants.CodeParamError, err)
	}

	if errortypes.IsDataNotFound(err) {
		return NewResultInfoWithCodeError(constants.CodeTargetNotFound, err)
	}

	if errortypes.IsPeerWait(err) {
		return NewResultInfoWithCodeError(constants.CodePeerWait, err)
	}

	if errortypes.IsPeerContinue(err) {
		return NewResultInfoWithCodeError(constants.CodePeerContinue, err)
	}

	if errortypes.IsURLNotReachable(err) {
		return NewResultInfoWithCodeError(constants.CodeURLNotReachable, err)
	}

	// IsConvertFailed
	return NewResultInfoWithCodeError(constants.CodeSystemError, err)
}

// NewResultInfoWithCodeError returns a new ResultInfo with code and error.
// And it will get the err.Error() as the value of ResultInfo.msg.
func NewResultInfoWithCodeError(code int, err error) ResultInfo {
	msg := err.Error()
	return NewResultInfoWithCodeMsg(code, msg)
}

// NewResultInfoWithCode returns a new ResultInfo with code
// and it will get the default msg corresponding to the code as the value of ResultInfo.msg.
func NewResultInfoWithCode(code int) ResultInfo {
	msg := constants.GetMsgByCode(code)
	return NewResultInfoWithCodeMsg(code, msg)
}

// NewResultInfoWithCodeMsg returns a new ResultInfo with code and specified msg.
func NewResultInfoWithCodeMsg(code int, msg string) ResultInfo {
	return NewResultInfo(code, msg, nil)
}

// NewResultInfoWithCodeData returns a new ResultInfo with code and specified data.
func NewResultInfoWithCodeData(code int, data interface{}) ResultInfo {
	return NewResultInfo(code, "", data)
}

// NewResultInfo returns a new ResultInfo.
func NewResultInfo(code int, msg string, data interface{}) ResultInfo {
	return ResultInfo{
		code: code,
		msg:  msg,
		data: data,
	}
}

func (r ResultInfo) Error() string {
	return fmt.Sprintf("{\"Code\":%d,\"Msg\":\"%s\"}", r.code, r.msg)
}

// SuccessCode return whether the code equals SuccessCode.
func (r ResultInfo) SuccessCode() bool {
	return r.code == constants.Success
}
