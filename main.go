package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lixin9311/authorize/storage"
	"github.com/lixin9311/osin"
)

const (
	secret = "secret"
)

type logger struct{}

func (l *logger) Printf(format string, v ...interface{}) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Output(3, fmt.Sprintf(format, v...))
}

func hello(c echo.Context) error {
	accessData := c.Get("AccessData").(*osin.AccessData)
	return c.String(200, fmt.Sprintf("hello %v", accessData))
}

func main() {
	exp := flag.Int("e", 3600, "expiration time")
	flag.Parse()
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.REFRESH_TOKEN}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	sconfig.AccessExpiration = int32(*exp)
	server := osin.NewServer(sconfig, storage.NewTestStorage())
	server.Logger = &logger{}
	h := &AutorizeHandler{server}
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.POST("/authorize", h.Authorize)
	e.POST("/token", h.Token)

	r := e.Group("/protected")
	r.Use(server.ValidatorMiddleware)
	r.GET("", hello)

	go e.Logger.Fatal(e.Start(":8080"))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Print("closing connection")
}
