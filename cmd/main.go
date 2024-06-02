package main

import (
	"context"
	"fmt"
	"go/test/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	hh := handlers.ResponseFunc()
	gb := handlers.ResponseFuncGoodbye()
	prod := handlers.GetProductsHandlerfunc()

	// custom server we created..

	sm := http.NewServeMux()

	sm.Handle("/goodbye", gb)

	sm.Handle("/products/", prod)
	sm.Handle("/", hh)

	srv := &http.Server{
		Addr:        ":8080",
		Handler:     sm,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
	}

	go func() {
		srv.ListenAndServe()
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// this <- is waiting for signal on the sigChan channel
	// once signal is recieved we print and shutdown with 30 seconds contect
	sig := <-sigChan
	fmt.Println("Received termination request. Closing the server gracefuly", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

}
