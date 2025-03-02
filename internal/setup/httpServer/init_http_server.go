package httpServer

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"quiz-app-be/internal/config"
	"time"

	"github.com/go-chi/chi"
)

type HttpServer struct {
	server    *http.Server  // сервер HTTP
	chControl chan struct{} // управление ожиданием завершения работы сервера
	lis       net.Listener
}

func Init(cfgHttp config.HttpServerConfig, cfg *config.Config, handlerFunc func(*config.Config) *chi.Mux) (comp *HttpServer, err error) {
	comp = &HttpServer{
		server: &http.Server{
			Addr:           fmt.Sprintf("%s:%d", cfgHttp.Host, cfgHttp.Port),
			Handler:        handlerFunc(cfg),
			ReadTimeout:    cfgHttp.ReadTimeout,
			WriteTimeout:   cfgHttp.WriteTimeout,
			IdleTimeout:    cfgHttp.IdleTimeout,
			MaxHeaderBytes: cfgHttp.MaxHeaderBytes,
		},
		chControl: make(chan struct{}),
	}

	if comp.lis, err = net.Listen("tcp", comp.server.Addr); err != nil {
		return
	}

	go func() {
		_ = comp.server.Serve(comp.lis)
		close(comp.chControl)
	}()

	return comp, nil
}

func (comp *HttpServer) Close() error {
	if comp.lis == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := comp.server.Shutdown(ctx); err != nil {
		if err := comp.lis.Close(); err != nil {
			return err
		}
	}

	<-comp.chControl
	return nil
}
