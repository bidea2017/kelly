// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package kelly

type dataContext interface {
	Set(interface{}, interface{}) dataContext
	Get(interface{}) interface{}
	MustGet(interface{}) interface{}
}
