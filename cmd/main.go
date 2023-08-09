package main

import (
	"os"

	"InterviewBackendSawitProGolang/generated"
	"InterviewBackendSawitProGolang/handler"
	"InterviewBackendSawitProGolang/pkg/middleware"
	"InterviewBackendSawitProGolang/repository"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"log"
	// "net/http"
)

func main() {
	e := echo.New()
	mw, err := middleware.NewMiddleware()
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	e.Use(echoMiddleware.Logger())
	e.Use(mw)
	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}
