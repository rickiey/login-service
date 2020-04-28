// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/9

package helper

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Response struct {
	IsSuccess bool        `json:"isSuccess"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

func responseOut(w http.ResponseWriter, msg string, status int, isSuccess bool, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errJson, _ := json.Marshal(Response{IsSuccess: isSuccess, Message: msg, Data: data})
	_, er := w.Write(errJson)
	if er != nil {
		logrus.Print(er.Error())
	}
}

func ResponseErr(w http.ResponseWriter, err error, status int) {
	responseOut(w, err.Error(), status, false, nil)
}

func ResponseData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	errJson, _ := json.Marshal(Response{IsSuccess: true, Data: data})
	_, er := w.Write(errJson)
	if er != nil {
		logrus.Print(er.Error())
	}
}

func ResponseMsg(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	errJson, _ := json.Marshal(Response{IsSuccess: true, Message: msg})
	_, er := w.Write(errJson)
	if er != nil {
		logrus.Print(er.Error())
	}
}
