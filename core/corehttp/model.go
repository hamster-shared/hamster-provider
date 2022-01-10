package corehttp

type Result struct {
	Code    uint64      `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func Success(result interface{}) Result {
	return Result{
		Code:    0,
		Type:    "success",
		Message: "",
		Result:  result,
	}
}

func BadRequest(message ...string) Result {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	}
	return Result{
		Code:    1,
		Type:    "fail",
		Message: msg,
	}
}
