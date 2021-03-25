package dataHandle

import "meetroom/err"

type Result struct {
	Data interface{} `json:"data"`

	Err     bool   `json:"err"`
	Message string `json:"mesage"`

	Code uint `json:"err_code"`
}

func OkResult(data interface{}) Result {
	return Result{
		Data:    data,
		Err:     false,
		Message: "",
		Code:    0,
	}
}

func FailureResult(e err.Exception) Result {
	return Result{
		Data:    nil,
		Err:     true,
		Message: e.Message + " | " + e.ExtraMessage,
		Code:    e.Code,
	}
}

func FailureFuncResult(eFunc func(string) err.Exception, extra string) Result {
	e := eFunc(extra)
	return FailureResult(e)
}
