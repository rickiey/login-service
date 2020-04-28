// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/2/2

package helper

import (
	"net/http"
)

// 依次从 Header 、url参数、 cookie 里获取 token, 有就返回
func GetToken(r *http.Request) string {

	var token string
	// get token form Header or URL or Cookies
	if t := r.Header.Get("token"); len(t) != 0 {
		token = t
	} else if t := r.URL.Query().Get("token"); len(t) != 0 {
		token = t
	} else if cookie, err := r.Cookie("access_token"); err == nil {
		token = cookie.Value
	}
	return token
}
