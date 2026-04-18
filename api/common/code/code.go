// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package code

type Code int

const Success Code = 10000

const (
	_ = iota + Success
	Error
	ParamErr
	Unauthorized
	Forbidden
	NotFound
	InternalServerError
)

var messages = map[Code]string{
	Success:             "success",
	Error:               "error",
	ParamErr:            "parameter error",
	Unauthorized:        "unauthorized",
	Forbidden:           "forbidden",
	NotFound:            "not found",
	InternalServerError: "internal server",
}

func (c Code) Message() string {
	msg, ok := messages[c]
	if !ok {
		return "Unknown"
	}

	return msg
}
