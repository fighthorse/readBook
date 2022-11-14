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
	router := routers.InitRouter()
	conf := setting.Config.Server
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.Port),
		Handler:        router,
		ReadTimeout:    conf.ReadTimeout * time.Second,
		WriteTimeout:   conf.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
