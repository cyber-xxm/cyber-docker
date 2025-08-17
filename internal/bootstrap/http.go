package bootstrap

import (
	"cyber-docker/internal/wirex"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func startHTTPServer(injector *wirex.Injector) error {

	e := gin.New()
	e.Use(gin.Recovery())
	injector.RegisterRouters(e)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      e,
		ReadTimeout:  time.Second * time.Duration(60),
		WriteTimeout: time.Second * time.Duration(60),
		IdleTimeout:  time.Second * time.Duration(10),
	}
	return srv.ListenAndServe()
}
