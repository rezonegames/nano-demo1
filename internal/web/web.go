package web

import (
	"context"
	"encoding/json"
	"github.com/lonng/nex"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tetris/config"
	"tetris/internal/web/api"
	"tetris/pkg/log"
)

func startupService() http.Handler {
	var (
		mux = http.NewServeMux()
	)

	nex.Before(logRequest)
	mux.Handle("/v1/user/", api.MakeLoginService())
	return AccessControl(OptionControl(mux))
}

func logRequest(ctx context.Context, r *http.Request) (context.Context, error) {
	if uri := r.RequestURI; uri != "/ping" {
		log.Debug("Method=%s, RemoteAddr=%s URL=%s", r.Method, r.RemoteAddr, uri)
	}
	return ctx, nil
}

func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		h.ServeHTTP(w, r)
	})
}

func OptionControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			json.NewEncoder(w).Encode(`{ "code": 0, "data": "success"}`)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func StartUp() {
	sc := config.ServerConfig
	go func() {
		// http service
		mux := startupService()
		log.Fatal(http.ListenAndServe(sc.ServerPort, mux))
	}()

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	// stop server
	select {
	case s := <-sg:
		log.Info("got signal: %s", s.String())
	}
}
