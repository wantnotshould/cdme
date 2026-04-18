// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package common

import (
	"encoding/json"
	"net/http"

	"code.cn/blog/api/common/code"
)

type response[T any] struct {
	Code    code.Code `json:"code"`
	Status  bool      `json:"status"`
	Message string    `json:"message"`
	Data    T         `json:"data,omitempty"`
}

func write[T any](w http.ResponseWriter, statusCode int, resp response[T]) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

func OK(w http.ResponseWriter) {
	write(w, http.StatusOK, response[any]{
		Code:    code.Success,
		Status:  true,
		Message: code.Success.Message(),
		Data:    nil,
	})
}

func OkData[T any](w http.ResponseWriter, data T) {
	write(w, http.StatusOK, response[T]{
		Code:    code.Success,
		Status:  true,
		Message: code.Success.Message(),
		Data:    data,
	})
}

func Fail(w http.ResponseWriter, code code.Code) {
	write(w, http.StatusOK, response[any]{
		Code:    code,
		Status:  false,
		Message: code.Message(),
		Data:    nil,
	})
}

func Custom(w http.ResponseWriter, status int, message string) {
	write(w, status, response[any]{
		Code:    code.Error,
		Status:  false,
		Message: message,
	})
}

func FailMessage(w http.ResponseWriter, message string) {
	write(w, http.StatusOK, response[any]{
		Code:    code.Error,
		Status:  false,
		Message: message,
		Data:    nil,
	})
}

func Unauthorized(w http.ResponseWriter) {
	write(w, http.StatusUnauthorized, response[any]{
		Code:    code.Unauthorized,
		Status:  false,
		Message: code.Unauthorized.Message(),
		Data:    nil,
	})
}

func Forbidden(w http.ResponseWriter) {
	write(w, http.StatusForbidden, response[any]{
		Code:    code.Forbidden,
		Status:  false,
		Message: code.Forbidden.Message(),
		Data:    nil,
	})
}

func NotFound(w http.ResponseWriter) {
	write(w, http.StatusNotFound, response[any]{
		Code:    code.NotFound,
		Status:  false,
		Message: code.NotFound.Message(),
		Data:    nil,
	})
}

func InternalServer(w http.ResponseWriter) {
	write(w, http.StatusNotFound, response[any]{
		Code:    code.InternalServerError,
		Status:  false,
		Message: code.InternalServerError.Message(),
		Data:    nil,
	})
}
