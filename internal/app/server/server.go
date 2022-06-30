package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/romik1505/fileserver/internal/app/config"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

type App struct {
	httpServer http.Server
}

func NewApp(ctx context.Context, handler http.Handler) *App {
	port := config.GetValue(config.Port)
	return &App{
		httpServer: http.Server{
			Addr:         port,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      handler,
		},
	}
}

func (a *App) Run() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func(ch chan os.Signal) {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Println(err.Error())
			done <- os.Interrupt
			return
		}
	}(done)

	tracer, closer := initJaeger("file-server")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	log.Printf("Server started on %s port", config.GetValue(config.Port))

	<-done
	defer close(done)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log.Println("Server gracefully closed")
	return a.httpServer.Shutdown(ctx)
}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &jaegerConfig.Configuration{
		ServiceName: service,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
