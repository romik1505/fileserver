package main

import (
	"context"
	"log"

	_ "github.com/gorilla/mux"
	"github.com/romik1505/fileserver/internal/app/handler"
	"github.com/romik1505/fileserver/internal/app/server"
	"github.com/romik1505/fileserver/internal/app/service"
)

func main() {
	ctx := context.Background()

	s := service.NewFileService()
	h := handler.NewHandler(s)

	app := server.NewApp(ctx, h.InitRoutes())
	if err := app.Run(); err != nil {
		log.Println(err.Error())
	}
}
