// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package kelly

import (
	"github.com/urfave/negroni"
	"net/http"
)

func wrapHandlerFunc(f HandlerFunc) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		f(newContext(rw, r, next))
	}
}