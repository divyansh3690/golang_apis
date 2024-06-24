package main

import (
	"context"
	"fmt"
	"go/test/data"
	"go/test/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	ph := handlers.GetProductsHandlerfunc()
	home := handlers.ResponseFunc()
	// gorilla mux router
	sm := mux.NewRouter()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading env file")
	}
	fmt.Print("env file is ", os.Getenv("GREETING"))
	connections := data.GetMongoDBFunctions()
	connections.Mongo_Connect()
	// we can create sub router with specific verbs like GET / POST with specific functions.
	getRoute := sm.Methods(http.MethodGet).Subrouter()
	getRoute.HandleFunc("/products", ph.GetReqProd)

	getRoute.HandleFunc(("/products/{id:[0-9]+}"), ph.GetProdByID)

	putRoute := sm.Methods(http.MethodPut).Subrouter()
	putRoute.HandleFunc("/products/{id:[0-9]+}", ph.UpdateProduct)
	putRoute.Use(ph.MiddlewaresHandlers)

	postRoute := sm.Methods(http.MethodPost).Subrouter()
	postRoute.HandleFunc("/products", ph.AddProduct)
	postRoute.Use(ph.MiddlewaresHandlers)

	getHomeRoute := sm.Methods(http.MethodGet).Subrouter()
	getHomeRoute.HandleFunc("/", home.ServeHTTP)

	delProdRoute := sm.Methods(http.MethodDelete).Subrouter()
	delProdRoute.HandleFunc("/products/{id:[0-9]+}", ph.RemoveProduct)

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
