package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fighthorse/readBook/common/setting"
	"github.com/fighthorse/readBook/common/validator"
	"github.com/fighthorse/readBook/routers"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	binding.Validator = new(validator.DefaultValidator)

	// user
	go UserServer()

	AdminServer()
}

func UserServer() {
	// User
	router := routers.InitUserRouter()
	conf := setting.Config.UserServer
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.Port),
		Handler:        router,
		ReadTimeout:    conf.ReadTimeout * time.Second,
		WriteTimeout:   conf.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	_ = s.ListenAndServe()
}

func AdminServer() {
	// ADMIN
	router := routers.InitRouter()
	conf := setting.Config.Server
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.Port),
		Handler:        router,
		ReadTimeout:    conf.ReadTimeout * time.Second,
		WriteTimeout:   conf.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	_ = s.ListenAndServe()
}
