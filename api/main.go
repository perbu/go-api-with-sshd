package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (a *API) Run(ctx context.Context, addr string) error {
	router := gin.Default()
	router.GET("/user/:name", a.getUser)
	router.POST("/user/:name/addpet", a.addPet)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		err := srv.Shutdown(ctxShutdown)
		if err != nil {
			log.Fatalln("Unexpected http shutdown error:", err)
		}
	}()
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}
