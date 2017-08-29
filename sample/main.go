// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/qjw/kelly"
	"github.com/qjw/kelly/middleware"
	"github.com/qjw/kelly/sessions"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

func initStore() sessions.Store {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       3,
	})
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatal("failed to connect redis")
	}

	store, err := sessions.NewRediStore(redisClient, []byte("abcdefg"))
	if err != nil {
		log.Print(err)
	}
	return store
}

func main() {
	store := initStore()

	router := kelly.New(
		middleware.Version("v1"),
		middleware.NoCache(),
		middleware.Secure(&middleware.SecureConfig{
			AllowedHosts: []string{"127.0.0.1:9090"},
		}),
		Middleware("v1", "v1", true),
	)

	router.SetNotFoundHandle(func(c *kelly.Context) {
		c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
	router.SetMethodNotAllowedHandle(func(c *kelly.Context) {
		c.WriteString(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	})

	InitParam(router)
	InitGroupMiddleware(router)
	InitRender(router)
	InitStatic(router)
	InitApiV1(router, store)
	InitCsrf(router)

	router.GET("/", func(c *kelly.Context) {
		c.WriteIndentedJson(http.StatusOK, kelly.H{
			"code":    "/",
			"value":   c.MustGet("v1"),                     // 获取context数据
			"session": c.GetDefaultCookie("session", "ss"), // 获取cookie数据
		})
	})

	router.Run(":9090")
}
