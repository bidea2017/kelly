// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package kelly

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"net/http"
)

type Router interface {
	GET(string, ...HandlerFunc)
	HEAD(string, ...HandlerFunc)
	OPTIONS(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	PUT(string, ...HandlerFunc)
	PATCH(string, ...HandlerFunc)
	DELETE(string, ...HandlerFunc)

	ServeHTTP(http.ResponseWriter, *http.Request)

	Group(string, ...HandlerFunc) Router

	// 设置404处理句柄
	SetNotFoundHandle(HandlerFunc)
	// 设置405处理句柄
	SetMethodNotAllowedHandle(HandlerFunc)
}

func (rt *router) wrapHandle(handles ...HandlerFunc) httprouter.Handle {
	var handle HandlerFunc = nil
	if len(rt.middlewares) > 0 {
		tmpHandle := negroni.New()
		for _, h := range rt.middlewares {
			tmpHandle.UseFunc(wrapHandlerFunc(h))
		}
		for _, v := range handles {
			tmpHandle.UseFunc(wrapHandlerFunc(v))
		}

		handle = func(c *Context) {
			tmpHandle.ServeHTTP(c, c.Request())
		}
	} else if len(handles) > 1 {
		// 没有中间件，但是有多个handle
		tmpHandle := negroni.New()
		for _, v := range handles {
			tmpHandle.UseFunc(wrapHandlerFunc(v))
		}

		handle = func(c *Context) {
			tmpHandle.ServeHTTP(c, c.Request())
		}
	} else {
		// 没有中间件，只有一个handle
		handle = handles[0]
	}
	return func(wr http.ResponseWriter, r *http.Request, params httprouter.Params) {
		r = mapContextFilter(wr, r, params)
		handle(newContext(wr, r, nil))
	}
}

type router struct {
	rt           *httprouter.Router
	path         string
	absolutePath string
	middlewares  []HandlerFunc
}

func (rt *router) SetNotFoundHandle(h HandlerFunc) {
	rt.rt.NotFound = h
}

func (rt *router) SetMethodNotAllowedHandle(h HandlerFunc) {
	rt.rt.MethodNotAllowed = h
}

func (rt *router) validatePath(path string, handles ...HandlerFunc) {
	if len(handles) < 1 {
		panic(fmt.Errorf("must have one handle"))
	}
	if len(path) < 1 {
		panic(fmt.Errorf("invalid path %s", path))
	}
	if path == "/" {
		return
	}
	if path[0] != '/' || path[len(path)-1] == '/' {
		panic(fmt.Errorf("invalid path %s,must beginwith (NOT endwith) /", path))
	}
}

func (rt *router) GET(path string, handles ...HandlerFunc) {
	rt.validatePath(path, handles...)
	rt.rt.GET(rt.absolutePath+path, rt.wrapHandle(handles...))
}

func (rt *router) HEAD(path string, handles ...HandlerFunc) {
	rt.validatePath(path, handles...)
	rt.rt.HEAD(rt.absolutePath+path, rt.wrapHandle(handles...))
}

func (rt *router) OPTIONS(path string, handles ...HandlerFunc) {
	rt.validatePath(path, handles...)
	rt.rt.OPTIONS(rt.absolutePath+path, rt.wrapHandle(handles...))
}

func (rt *router) POST(path string, handles ...HandlerFunc) {
	rt.validatePath(path, handles...)
	rt.rt.POST(rt.absolutePath+path, rt.wrapHandle(handles...))
}

func (rt *router) PUT(path string, handles ...HandlerFunc) {
	rt.validatePath(path, handles...)
	rt.rt.PUT(rt.absolutePath+path, rt.wrapHandle(handles...))
}

func (rt *router) PATCH(path string, handles ...HandlerFunc) {
	rt.validatePath(path, handles...)
	rt.rt.PATCH(rt.absolutePath+path, rt.wrapHandle(handles...))
}

func (rt *router) DELETE(path string, handles ...HandlerFunc) {
	rt.validatePath(path)
	rt.rt.DELETE(rt.absolutePath+path, rt.wrapHandle(handles...))
}

func (rt *router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rt.rt.ServeHTTP(rw, r)
}

func (r *router) Group(path string, handlers ...HandlerFunc) Router {
	tmpHandles := r.middlewares
	if len(tmpHandles) == 0 {
		tmpHandles = handlers
	} else if len(handlers) == 0 {
		tmpHandles = make([]HandlerFunc, len(r.middlewares))
		copy(tmpHandles, r.middlewares)
	} else {
		tmpHandles = make([]HandlerFunc, len(r.middlewares)+len(handlers))
		copy(tmpHandles, r.middlewares)
		for i, v := range handlers {
			tmpHandles[len(r.middlewares)+i] = v
		}
	}

	rt := &router{
		rt:           r.rt,
		path:         path,
		absolutePath: r.absolutePath + path,
		middlewares:  tmpHandles,
	}
	return rt
}

func newRouterImp(handlers ...HandlerFunc) *router {
	httpRt := httprouter.New()
	rt := &router{
		rt:           httpRt,
		path:         "",
		absolutePath: "",
	}

	if len(handlers) > 0 {
		rt.middlewares = make([]HandlerFunc, len(handlers))
		copy(rt.middlewares, handlers)
	}

	return rt
}

func newRouter(handlers ...HandlerFunc) Router {
	return newRouterImp(handlers...)
}
