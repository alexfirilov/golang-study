package main

import (
	"context"
	"fmt"
	"golang-study/internal/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fmt.Println("Starting Netscribe server...")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	http.HandleFunc("/health", api.HealthCheck)
	http.HandleFunc("/servers", api.handleServers)
	http.HandleFunc("/documents", api.handleDocuments)
	go api.InfrastructureWorker()

	go func() { srv.ListenAndServe() }()

	fmt.Println("Server started successfully!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutdown signal received. Gracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}
	fmt.Println("Waiting for background A agents to finish...")
	wg.Wait()
	fmt.Println("Server exiting")
}
