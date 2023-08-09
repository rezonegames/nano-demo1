package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tetris/config"
	"tetris/internal/web/api"
	"tetris/pkg/log"
)

func Register(r *gin.Engine) {
	/**
	 */
	r.POST("/v1/login", api.QueryHandler)
}

func StartUp() {
	sc := config.ServerConfig

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"x-checksum"},
	}))
	Register(r)
	srv := &http.Server{
		Addr:    sc.ServerPort,
		Handler: r,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return
		}
	}()
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	// stop server
	select {
	case s := <-sg:
		log.Info("got signal: %s", s.String())
	}
}
