// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/11

package helper

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"login-service/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

const FormatErrorMsg = "验证数据失败(%s)"

func Bind(r *http.Request, dto interface{}) error {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logrus.Error("ReadAll data: ", string(b))
		return err
	}
	// Unmarshal
	err = json.Unmarshal(b, dto)
	if err != nil {
		logrus.Error("Unmarshal again error,data: ", string(b), err)
		return err
	}
	return err
}

func ValidateError(errs error) string {
	var invalidParamsMessage []string
	if ve, ok := errs.(validator.ValidationErrors); ok {
		for _, err := range ve {
			//generate validation error message
			invalidParamsMessage = append(invalidParamsMessage, utils.ToSnakeCase(err.Field()))
		}
		return fmt.Sprintf(FormatErrorMsg, strings.Join(invalidParamsMessage, ","))
	}
	return "验证数据失败"
}
